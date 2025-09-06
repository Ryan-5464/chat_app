function PulseElement(elem) {
    elem.classList.remove(APP.CLS.GEN.PULSE)
    void elem.offsetWidth; 
    elem.classList.add(APP.CLS.GEN.PULSE)
}
