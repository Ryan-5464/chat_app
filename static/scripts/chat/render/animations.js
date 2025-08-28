function PulseElement(elem) {
    elem.classList.remove(APP.CLS.PULSE)
    void elem.offsetWidth; 
    elem.classList.add(APP.CLS.PULSE)
}
