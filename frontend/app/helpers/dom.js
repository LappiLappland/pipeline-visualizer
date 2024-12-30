export function isBottomVisible(element, margin = 0) {
    const rect = element.getBoundingClientRect();
    const viewportHeight = window.innerHeight || document.documentElement.clientHeight;

    return rect.bottom >= -margin && rect.bottom <= viewportHeight + margin;
}

export function scrollToBottomOnBelowView(element, containerHeight, extra = 0) {
    const rect = element.getBoundingClientRect();
    const viewportHeight = window.innerHeight || document.documentElement.clientHeight;

    const isInLowerHalf = rect.bottom > viewportHeight / 2;

    if (isInLowerHalf) {
        element.scrollTo({
            top: containerHeight + extra,
            behavior: 'instant',
        });
    }
}