/*
saveToken stores a token in local storage with an expiration time of 24 hours.

Parameters:
  - token (string): The token to be saved.

Returns:
  - None
*/
export const saveToken = (token) => {
    const expireAt = new Date().getTime() + 24 * 60 * 60 * 1000; // 24 hours
    const tokenData = {
        value: token,
        expireAt: expireAt
    };
    localStorage.setItem('token', JSON.stringify(tokenData));
};

/*
getToken retrieves a token from local storage if it is still valid.

Parameters:
  - None

Returns:
  - (string|null): The token value if valid, or null if expired or not found.
*/
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

/*
isAuthenticated checks whether a valid token exists in local storage.

Parameters:
  - None

Returns:
  - (boolean): True if a valid token exists, false otherwise.
*/
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
