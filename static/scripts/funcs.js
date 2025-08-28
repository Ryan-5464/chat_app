function RemoveAllListeners(elem) {
    const newElem = elem.cloneNode(true)
    elem.replaceWith(newElem)
    return newElem
}

function GetElemByDataTag(tag, val) {
    return document.querySelector(`[data-${tag}="${val}"]`);
};

function GetDataAttribute(elem, tag) {
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

function DeleteElementByDataTag(tagName, tagVal) {
    idStr = `[data-${tagName}="${tagVal}"]`
    const elem = document.querySelector(idStr);
    if (elem) {
        elem.remove();
    } else {
        throw new Error(`Failed to find element for identifier = ${idStr}`);
    };
};

function GetClosestTargetByData(e, target) {
    return e.target.closest(`[data-${target}]`)
}