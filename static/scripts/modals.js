window.addEventListener("click", function (e) {
    QSelectAllByClass(document, APP.CLS.MODAL.TAG).forEach(modal => {
        const modalContent = modal.querySelector(APP.CLS.MODAL.CONTENT);
        if (modal.classList.contains(APP.CLS.GEN.OPEN) && !modalContent.contains(e.target)) {
            modal.__controller.Close();
        };
    });
});

function replaceWithInput(elem, placeholder, id) {
    const input = CreateElement({elemType:'input', id:id, classes:[APP.CLS.GEN.INPUT_ELEM]});
    input.__oldtext = elem.textContent;
    input.type = 'text';
    input.name = 'Name';
    input.value = placeholder;
    input.required = true;
    elem.appendChild(input);
    return input;
};

function replaceWithTextArea(elem, placeholder, id) {
    const textarea = CreateElement({elemType:'textarea', id:id, classes:[APP.CLS.GEN.INPUT_ELEM, APP.CLS.GEN.TEXT_AREA]});
    textarea.name = 'Name';
    textarea.value = placeholder;
    textarea.required = true;
    elem.appendChild(textarea);
    configureTextArea(textarea); 
    return textarea;
};

function configureTextArea(textarea) {
    function resizeTextarea() {
        textarea.style.height = 'auto';
        textarea.style.height = textarea.scrollHeight + 'px';
    }
    textarea.addEventListener('input', resizeTextarea);
    resizeTextarea(); 
};