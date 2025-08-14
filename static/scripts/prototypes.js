Function.prototype.next = function (nextFn) {
  const prevFn = this;

  return function (...args) {
    const result = prevFn(...args);
    return nextFn(result);
  };
};
