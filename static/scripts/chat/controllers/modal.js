function CloseModal (modal) {
    const modalContent = QSelectByClass(modal, APP.CLS.MODAL.CONTENT)
    modalContent.classList.remove(APP.CLS.GEN.OPENING);
    modalContent.classList.add(APP.CLS.GEN.CLOSING);

    setTimeout(() => {
        modal.classList.remove(APP.CLS.GEN.OPEN);
        modalContent.classList.remove(APP.CLS.GEN.CLOSING);
    }, MODAL_CLOSE_DELAY);
}

function OpenModalAt(modal, clientX, clientY) {
    const modalContent = QSelectByClass(modal, APP.CLS.MODAL.CONTENT);
    document.querySelectorAll(APP.CLS.MODAL.OPEN).forEach(openModal => {
        if (openModal !== modal && openModal.__controller) {
            openModal.__controller.Close();
        }
    });

    const padding = 10;
    const maxLeft = window.innerWidth - modalContent.offsetWidth - padding;
    const maxTop = window.innerHeight - modalContent.offsetHeight - padding;

    modalContent.style.left = Math.min(clientX, maxLeft) + "px";
    modalContent.style.top = Math.min(clientY, maxTop) + "px";

    modalContent.classList.remove(APP.CLS.GEN.OPENING, APP.CLS.GEN.CLOSING);
    void modalContent.offsetWidth; 

    modal.classList.add(APP.CLS.GEN.OPEN);
    modalContent.classList.add(APP.CLS.GEN.OPENING);
};





