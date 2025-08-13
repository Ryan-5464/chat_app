
function BuildURL(endpoint) {
    return `${BASEURL}/${endpoint}`  
}

function BuildURLWithParams(endpoint, params) {
    const query = new URLSearchParams(params).toString();
    return `${BASEURL}/${endpoint}?${query}`;
}
