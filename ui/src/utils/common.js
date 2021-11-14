const unit = ["B", "KB", "MB", "TB"]

function formatSize(size) {
    let idx = 0;
    while (size > 1024) {
        size /= 1024;
        idx++;
    }
    return size.toFixed(2) + unit[idx];
}


function toInt(v, def = null) {
    if (!v) {
        return def;
    }
    try {
        const i = parseInt(v);
        if (isNaN(i)) {
            return def;
        } else {
            return i;
        }
    } catch (e) {
        return def;
    }
}

export {formatSize,toInt}