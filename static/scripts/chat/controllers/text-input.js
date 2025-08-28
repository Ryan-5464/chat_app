
function replaceTextInputWithText(openInput, textDivId) {
    if (!openInput) { return; }
    const openInputText = CreateElement({classes:[textDivId], innerHTML:openInput.__oldtext});
    openInput.replaceWith(openInputText);
};

function replaceTextWithTextInput(container, textDivId, inputId, isMsg) {
    const elemText = QSelectByClass(container, textDivId);
    let input
    if (isMsg) {
        input = replaceWithTextArea(elemText, elemText.textContent, inputId);
    } else {
        input = replaceWithInput(elemText, elemText.textContent, inputId);
    };
    input.focus();
    return input;
};

function textInputController(container, submitTextHandler, inputId, textDivId, isMsg=false) {
    const openInput = document.getElementById(inputId)
    replaceTextInputWithText(openInput, textDivId) 
    if (openInput) { openInput.remove()}

    const input = replaceTextWithTextInput(container, textDivId, inputId, isMsg);
    input.addEventListener('keydown', (e) => {
        if (e.key === "Enter") {
            submitTextHandler(input.value);
        };
    });

    const outsideClickHandler = function (event) {
        if (event.target !== input && !input.contains(event.target)) {
            replaceTextInputWithText(input, textDivId);
            document.removeEventListener('click', outsideClickHandler);
        }
    };

    // Delay attaching the event to allow the current stack to complete
    setTimeout(() => {
        document.addEventListener('click', outsideClickHandler);
    }, 0);
}
