import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { uploadPDF, fetchUserInfo } from '../api';
import { logout} from "../utils";

const Translate = () => {
    const [file, setFile] = useState(null);
    const [lang, setLang] = useState('eng'); // 选择语言
    const [downloadLink, setDownloadLink] = useState(''); // 服务端返回的文件链接
    const [isLoading, setIsLoading] = useState(false); // 是否加载中
    const [isSidebarVisible, setIsSidebarVisible] = useState(false); // 控制侧边栏显示
    const [userInfo, setUserInfo] = useState({ username: '', balance: 0 }); // 用户信息
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        alert('Logout successfully');
        navigate(0);
    };

    const handleUpload = async (e) => {
        e.preventDefault();
        if (!file) {
            alert('Please upload a file');
            return;
        }

        const formData = new FormData();
        formData.append('document', file);
        formData.append('lang', lang);

        setIsLoading(true);

        try {
            const { data } = await uploadPDF(formData);
            setDownloadLink(data.data);
        } catch (error) {
            if (error.response?.status === 401) {
                alert('Unauthorized, please sign in');
                navigate('/login');
            } else {
                alert('Upload failed, please try again later');
            }
        } finally {
            setIsLoading(false);
        }
    };

    const toggleSidebar = async () => {
        setIsSidebarVisible(!isSidebarVisible);
        if (!isSidebarVisible) {
            try {
                const response = await fetchUserInfo();
                setUserInfo({username: response.data.username, balance: response.data.balance});
            } catch (error) {
                console.error('Failed to get user information', error);
            }
        }
    };

    return (
        <div style={{position: 'relative'}}>
            <form onSubmit={handleUpload}>
                <h2>PDF Translate</h2>
                <div className="custom-file-upload">
                    <input
                        id="file-upload"
                        type="file"
                        accept=".pdf"
                        onChange={(e) => setFile(e.target.files[0])}
                        disabled={isLoading}
                    />
                    <label htmlFor="file-upload">
                        {file ? `Selected: ${file.name}` : "Click to select file"}
                    </label>
                </div>
                <p>Source Language</p>
                <select
                    value={lang}
                    onChange={(e) => setLang(e.target.value)}
                    disabled={isLoading}
                >
                    <option value="eng">English</option>
                    <option value="ara">Arabic</option>
                    <option value="fra">French</option>
                    <option value="rus">Russian</option>
                    <option value="spa">Spanish</option>
                </select>
                <button type="submit" disabled={isLoading}>
                    {isLoading ? 'Processing...' : 'Submit'}
                </button>
                {isLoading && <p>Processing, please wait</p>}
                {downloadLink && (
                    <div>
                        <p>翻译完成！</p>
                        <a href={downloadLink} target="_blank" rel="noopener noreferrer">下载翻译文件</a>
                    </div>
                )}
            </form>

            <div style={{marginTop: '10px', display: 'flex', alignItems: 'center', gap: '10px'}}>
                <button onClick={handleLogout}>Logout</button>
                <button className="user-info-button" onClick={toggleSidebar}>
                    User Info
                </button>
            </div>
            {/* 侧边栏 */}
            {isSidebarVisible && (
                <div>
                    <h3>User Info</h3>
                    <p>Username: {userInfo.username}</p>
                    <p>Remaining Page Credit: {userInfo.balance}</p>
                </div>
            )}
        </div>
    );
};

export default Translate;
