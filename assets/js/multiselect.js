function categoriesClicked(multiselectId) {
    const multiSelect = document.getElementById(multiselectId);
    multiSelect.classList.toggle("open");
}

function categoriyListChanged(optionsListId, selectBoxId) {
    const optionsList = document.getElementById(optionsListId);
    const selectedCategories = Array.from(
                        optionsList.querySelectorAll("input:checked")
                    ).map((input) => input.value);

    const selectedCategoriesTexts = Array.from(
                            optionsList.querySelectorAll("input:checked")
                        ).map((input) => input.parentElement.textContent);

    const selectBox = document.getElementById(selectBoxId);
    selectBox.textContent = selectedCategories.length
        ? selectedCategoriesTexts.join(", ")
        : "Select options";
    
}