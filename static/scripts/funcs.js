function RemoveAllListeners(elem) {
    const newElem = elem.cloneNode(true)
    elem.replaceWith(newElem)
    return newElem
}

function GetElemByDataTag(elem, tag, val) {
    return elem.querySelector(`[data-${tag}="${val}"]`);
};

function GetDataAttribute(elem, tag) {
        if (!elem) { return null; }
    return elem.dataset[tag];
};

function QSelectByClass(elem, cls) {
    return elem.querySelector(".".concat(cls));
};

function QSelectAllByClass(elem, cls) {
    return elem.querySelectorAll(".".concat(cls));
};

function QSelectById(elem, id) {
    return elem.querySelector("#".concat(id));
};

function QSelectAllById(elem, id) {
    return elem.querySelectorAll("#".concat(id));
};

function DeleteElementByDataTag(elem, tagName, tagVal) {
    idStr = `[data-${tagName}="${tagVal}"]`
    const el = elem.querySelector(idStr);
    if (el) {
        el.remove();
    } else {
        throw new Error(`Failed to find element for identifier = ${idStr}`);
    };
};

function GetClosestTargetByData(e, target) {
    return e.target.closest(`[data-${target}]`)
}