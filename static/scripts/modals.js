window.addEventListener("click", function (e) {
    document.querySelectorAll(".modal").forEach(modal => {
        const modalContent = modal.querySelector(".modal-content");
        if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
            modal.__controller.Close()
        }
    });
});

function replaceWithInput(elem, placeholder, id) {
    const input = document.createElement('input');
    input.__oldtext = elem.textContent;
    elem.innerHTML = '';
    input.id = id;
    input.type = 'text';
    input.name = 'Name';
    input.value = placeholder;
    input.required = true;
    input.className = 'input-elem';
    elem.appendChild(input);
    return input;
}

function replaceWithTextArea(elem, placeholder, id) {
    const textarea = document.createElement('textarea');
    textarea.id = id;
    textarea.name = 'Name';
    textarea.value = placeholder;
    textarea.required = true;
    textarea.classList.add('input-elem', 'textArea');

    elem.innerHTML = '';
    elem.appendChild(textarea);

    configureTextArea(textarea); // Pass the element directly

    return textarea;
}

function configureTextArea(textarea) {
    function resizeTextarea() {
        textarea.style.height = 'auto';
        textarea.style.height = textarea.scrollHeight + 'px';
    }

    textarea.addEventListener('input', resizeTextarea);
    resizeTextarea(); // Resize on load
}