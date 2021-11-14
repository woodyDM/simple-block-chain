import {ErrorBox} from "./pop";
import apis from './apis';

//fn should return Promise<V>
const Wrap = function (fn) {
    return function () {
        return new Promise((resolve, reject) => {
            fn().then(res => {
                if (res.status === 200 && res.data.code === 0) {
                    resolve(res.data.data);
                } else {
                    const message = res.data ? res.data.msg : "未知错误";
                    reject(new Error(message));
                }
            }).catch(e => reject(e));
        });
    }
}
/**
 *
 * @param fn promise
 * @param fe func(e)
 * @returns {function(): Promise<unknown>}
 */
const wrapExp = function (fn, fe) {
    return function () {
        return new Promise((resolve, reject) => {
            Wrap(fn)().then(d => {
                resolve(d);
            }).catch(e => fe(e))
        });
    }
}

const WrapIgnore = function (fn) {
    return wrapExp(fn, console.log)
};

const WrapX = function (fn) {
    return wrapExp(fn, e => ErrorBox(e.message));
}

const WrapAdmin = function (history, fn) {
    return function () {
        return new Promise((resolve, reject) => {
            Wrap(fn)().then(resolve)
                .catch(e => {
                    if (e.response.status === 403) {
                        history.push(apis.page.login);
                    } else {
                        ErrorBox(e.message);
                    }
                })
        });
    }
}

export {WrapX, WrapAdmin, WrapIgnore}