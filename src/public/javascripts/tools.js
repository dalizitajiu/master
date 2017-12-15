const HOST = "http://127.0.0.1:8080";
const EMAIL_REG = new RegExp("^[a-z0-9]+([._\\-]*[a-z0-9])*@([a-z0-9]+[-a-z0-9]*[a-z0-9]+.){1,63}[a-z0-9]+$");

const getAticleToken = async function(articleId) {
    if (!Number.isInteger(articleId)) {
        return null
    }
    const url = `/article/gettoken`;
    const response = await axios.get(url);
    return response.data;
}
const getArticleByid = async function(id) {
    if (!Number.isInteger(id) || id <= 0) {
        console.log("wrong args");
        return null;
    }

    const url = `/article/${id}`;
    const response = await axios.get(url);
    return response.data

}
const getArticleList = async function() {
    const url = `/article/abstractlist`;
    const response = await axios.get(url);
    return response.data.data;
}
const getArticleByRid = async function() {
    const url = `/article/getones`;
    const response = await axios.get(url);
    return response.data;
}
const checkEmail = function(email) {
    return EMAIL_REG.test(email);
}
const login = async function(email, pwd) {
    if (!checkEmail(email)) {
        console.log("wrong args");
        alert("参数错误");
        return;
    }
    const url = `/user/login`
    const param = new URLSearchParams()
    param.append("email", email)
    param.append("pwd", pwd)
    const response = await axios.post(url, param);
    return response.data;
}