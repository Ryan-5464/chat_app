function RemoveAllListeners(elem) {
    const newElem = elem.cloneNode(true)
    elem.replaceWith(newElem)
    return newElem
}
