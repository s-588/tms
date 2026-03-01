// JavaScript helper functions for table filters and operations
function toggleAllCheckboxes(source, className) {
    const checkboxes = document.getElementsByClassName(className);
    for (let i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = source.checked;
    }
}

function getSelectedIds(className) {
    const checkboxes = document.getElementsByClassName(className);
    const selectedIds = [];
    for (let i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            selectedIds.push(checkboxes[i].value);
        }
    }
    return selectedIds;
}

// Initialize date pickers and other filter components
document.addEventListener('DOMContentLoaded', function() {
    // Initialize tooltips
    const tooltips = document.querySelectorAll('[data-tooltip]');
    tooltips.forEach(tooltip => {
        tooltip.addEventListener('mouseenter', function() {
            const tooltipText = this.getAttribute('data-tooltip');
            // Create tooltip element
            const tooltipEl = document.createElement('div');
            tooltipEl.className = 'tooltip is-tooltip-multiline';
            tooltipEl.textContent = tooltipText;
            tooltipEl.style.position = 'absolute';
            tooltipEl.style.zIndex = '9999';
            document.body.appendChild(tooltipEl);
            
            // Position tooltip
            const rect = this.getBoundingClientRect();
            tooltipEl.style.top = (rect.top - tooltipEl.offsetHeight - 10) + 'px';
            tooltipEl.style.left = (rect.left + rect.width / 2 - tooltipEl.offsetWidth / 2) + 'px';
            
            this.tooltipElement = tooltipEl;
        });
        
        tooltip.addEventListener('mouseleave', function() {
            if (this.tooltipElement) {
                this.tooltipElement.remove();
                this.tooltipElement = null;
            }
        });
    });
    
    // Initialize range filter validation
    const rangeInputs = document.querySelectorAll('input[type="number"][name$="_min"], input[type="number"][name$="_max"]');
    rangeInputs.forEach(input => {
        input.addEventListener('change', function() {
            const name = this.getAttribute('name');
            const isMin = name.endsWith('_min');
            const baseName = isMin ? name.replace('_min', '') : name.replace('_max', '');
            const minInput = document.querySelector(`input[name="${baseName}_min"]`);
            const maxInput = document.querySelector(`input[name="${baseName}_max"]`);
            
            if (minInput.value && maxInput.value && parseFloat(minInput.value) > parseFloat(maxInput.value)) {
                this.setCustomValidity('Minimum value cannot be greater than maximum value');
                this.reportValidity();
            } else {
                this.setCustomValidity('');
            }
        });
    });
});
