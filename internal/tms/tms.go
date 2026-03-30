package tms

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/s-588/tms/internal/db"
	"github.com/shopspring/decimal"
)

const (
	M = 0.3 // max penalty (30%)
	k = 3.0
)

var (
	fuelPricePerLiter = decimal.NewFromFloat(2.99)
)

func CalculateClientDiscount(ctx context.Context, clientID int32, db db.DB) (float64, error) {
	total, canceled, err := db.CountClientsOrders(ctx, clientID)
	if err != nil {
		return 0, fmt.Errorf("count client orders: %w", err)
	}

	if total == 0 {
		// First order → 20% discount
		return 0.20, nil
	}

	// Avoid division by zero if totalOrders == 0 (already handled)
	n := float64(canceled)
	T := float64(total)

	// Penalty formula: M * (1 - e^(-k*n)) * (n/T)^0.7
	penalty := M * (1 - math.Exp(-k*n)) * math.Pow(n/T, 0.7)

	// Initial discount 20% minus penalty, cannot go below 0
	discount := 0.20 - penalty
	if discount < 0 {
		discount = 0
	}
	return discount, nil
}

type CalculateOrderCostArgs struct {
	ClientID                         int32
	PriceID                          int32
	Weight                           int64
	FuelConsumption, PayloadCapacity int32
	NodeStartID, NodeEndID           int32
}

func CalculateOrderCost(ctx context.Context, db db.DB, args CalculateOrderCostArgs) (decimal.Decimal, error) {
	slog.Debug("🚀 CalculateOrderCost started",
		"client_id", args.ClientID,
		"weight_kg", args.Weight,
		"payload_kg", args.PayloadCapacity,
		"distance_m", args.NodeStartID, // если нужно, замени на реальное
	)

	// 1. Расстояние (h в км)
	distanceMeters, err := db.CalculateDistance(ctx, args.NodeStartID, args.NodeEndID)
	if err != nil {
		slog.Error("❌ failed to calculate distance", slog.Any("error", err))
		return decimal.Decimal{}, fmt.Errorf("calculate distance: %w", err)
	}
	h := decimal.NewFromFloat(distanceMeters)
	slog.Debug("📏 distance", slog.String("h_km", h.String()))

	// 2. Конверсия в тонны
	wKg := decimal.NewFromInt(int64(args.Weight))
	vKg := decimal.NewFromInt(int64(args.PayloadCapacity))
	w := wKg.Div(decimal.NewFromInt(1000))
	v := vKg.Div(decimal.NewFromInt(1000))
	slog.Debug("⚖️ weight and payload (tons)",
		slog.String("w_t", w.String()),
		slog.String("v_t", v.String()),
	)

	if w.GreaterThan(v) {
		slog.Error("⛔ weight exceeds payload", slog.String("w", w.String()), slog.String("v", v.String()))
		return decimal.Decimal{}, fmt.Errorf("weight (%s t) exceeds payload (%s t)", w, v)
	}

	// 3. Стоимость топлива
	r := decimal.NewFromInt32(args.FuelConsumption)
	c := fuelPricePerLiter // ← замени на реальную переменную/аргумент

	base := r.Mul(decimal.NewFromFloat(0.4)).Mul(w)
	fuelLiters := h.Div(decimal.NewFromInt(100)).Mul(base)
	fuelCost := fuelLiters.Mul(c)

	slog.Debug("⛽ fuel calculation",
		slog.String("r", r.String()),
		slog.String("base_l_per_100km", base.String()),
		slog.String("fuel_liters", fuelLiters.String()),
		slog.String("fuel_cost_before_factor", fuelCost.String()),
	)

	// 4. Коэффициент загрузки
	loadRatio := w.Div(v)
	factor := decimal.NewFromInt(1).
		Sub(loadRatio).
		Div(decimal.NewFromInt(2)).
		Add(decimal.NewFromInt(1))

	C := fuelCost.Mul(factor)
	slog.Debug("📈 load factor",
		slog.String("load_ratio", loadRatio.String()),
		slog.String("factor", factor.String()),
		slog.String("C_fuel_final", C.String()),
	)
	price, err := db.GetPriceByID(ctx, args.PriceID)
	// 4. Надбавки: C * Kh * Kw * Kx
	Kh := price.Distance // коэффициент увеличения стоимости за расстояние
	Kw := price.Weight   // коэффициент увеличения стоимости за вес груза
	totalBeforeDiscount := C.Mul(Kh).Mul(Kw)
	slog.Debug("📊 coefficients and subtotal",
		slog.String("Kh", Kh.String()),
		slog.String("Kw", Kw.String()),
		slog.String("total_before_discount", totalBeforeDiscount.String()),
	)

	// 6. Скидка клиента
	discount, err := CalculateClientDiscount(ctx, args.ClientID, db)
	if err != nil {
		slog.Error("❌ failed to calculate discount", slog.Any("error", err))
		return decimal.Decimal{}, fmt.Errorf("calculate client discount: %w", err)
	}
	discountFactor := decimal.NewFromFloat(1 - discount)
	slog.Debug("🎟️ client discount",
		slog.Float64("discount_percent", discount*100),
		slog.String("discount_factor", discountFactor.String()),
	)

	// 7. Итоговая цена
	totalPrice := totalBeforeDiscount.Mul(discountFactor)
	slog.Info("✅ order cost calculated",
		slog.String("total_price", totalPrice.String()),
		slog.Float64("total_price_float", totalPrice.InexactFloat64()),
	)

	return totalPrice, nil
}
