window.addEventListener("click", function (e) {
    document.querySelectorAll(".modal").forEach(modal => {
        const modalContent = modal.querySelector(".modal-content");
        if (modal.classList.contains("open") && !modalContent.contains(e.target)) {
            modal.__controller.Close()
        }
    });
});

function replaceWithInput(elem, placeholder) {
    const input = document.createElement('input');
    input.__oldtext = elem.innerHTML
    elem.innerHTML = ''
    input.id = 'chat-name-input'
    input.type = 'text';
    input.name = 'Name';
    input.placeholder = placeholder;
    input.required = true;
    input.className = 'input-elem';
    elem.appendChild(input);
    return input
}


