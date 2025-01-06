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
        alert('成功登出');
        navigate(0);
    };

    const handleUpload = async (e) => {
        e.preventDefault();
        if (!file) {
            alert('请选择一个文件');
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
                alert('未授权，请登录');
                navigate('/login');
            } else {
                alert('上传失败，请稍后重试');
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
                console.error('获取用户信息失败', error);
            }
        }
    };

    return (
        <div style={{position: 'relative'}}>
            <form onSubmit={handleUpload}>
                <h2>PDF 翻译</h2>
                <div className="custom-file-upload">
                    <input
                        id="file-upload"
                        type="file"
                        accept=".pdf"
                        onChange={(e) => setFile(e.target.files[0])}
                        disabled={isLoading}
                    />
                    <label htmlFor="file-upload">
                        {file ? `已选择: ${file.name}` : "点击选择文件"}
                    </label>
                </div>
                <select
                    value={lang}
                    onChange={(e) => setLang(e.target.value)}
                    disabled={isLoading}
                >
                    <option value="eng">英文</option>
                    <option value="ara">阿拉伯文</option>
                    <option value="fra">法文</option>
                    <option value="rus">俄文</option>
                    <option value="spa">西班牙文</option>
                </select>
                <button type="submit" disabled={isLoading}>
                    {isLoading ? '处理中...' : '提交'}
                </button>
                {isLoading && <p>处理中，请稍候...</p>}
                {downloadLink && (
                    <div>
                        <p>翻译完成！</p>
                        <a href={downloadLink} target="_blank" rel="noopener noreferrer">下载翻译文件</a>
                    </div>
                )}
            </form>

            <div style={{marginTop: '10px', display: 'flex', alignItems: 'center', gap: '10px'}}>
                <button onClick={handleLogout}>登出</button>
                <button className="user-info-button" onClick={toggleSidebar}>
                    用户信息
                </button>
            </div>
            {/* 侧边栏 */}
            {isSidebarVisible && (
                <div>
                    <h3>用户信息</h3>
                    <p>用户名: {userInfo.username}</p>
                    <p>剩余额度: {userInfo.balance}</p>
                    <button onClick={toggleSidebar}>
                        关闭
                    </button>
                </div>
            )}
        </div>
    );
};

export default Translate;
