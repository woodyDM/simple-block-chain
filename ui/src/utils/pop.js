import {notification} from "antd";

const ErrorBox = function (msg) {
    const ms = (msg ? msg : "未知错误")
    notification.error({
        message: '错误',
        description: ms
    });
}

const OkBox = function (msg) {
    const ms = (msg ? msg : "完成")
    notification.success({
        message: '成功',
        description: ms
    });
}

export {ErrorBox, OkBox}
