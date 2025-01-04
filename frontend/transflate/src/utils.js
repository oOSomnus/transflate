export const saveToken = (token) => {
    const expireAt = new Date().getTime() + 24 * 60 * 60 * 1000; // 24 hours
    const tokenData = {
        value: token,
        expireAt: expireAt
    };
    localStorage.setItem('token', JSON.stringify(tokenData));
};

export const getToken = () => {
    const tokenData = JSON.parse(localStorage.getItem('token')); // fetch and parse
    if (tokenData) {
        const currentTime = new Date().getTime();
        if (currentTime > tokenData.expireAt) {
            localStorage.removeItem('token'); // remove token
            return null;
        }
        return tokenData.value; // return token value
    }
    return null; // token not exists
};

export const isAuthenticated = () => {
    const tokenData = JSON.parse(localStorage.getItem('token'));
    if (tokenData) {
        const currentTime = new Date().getTime();
        if (currentTime > tokenData.expireAt) {
            localStorage.removeItem('token');
            return false;
        }
        return true; // token still valid
    }
    return false; // token not exists
};
export const logout = () => {
    localStorage.removeItem('token');
};
