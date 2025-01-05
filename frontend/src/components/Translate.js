import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { uploadPDF } from '../api';

const Translate = () => {
    const [file, setFile] = useState(null);
    const [lang, setLang] = useState('eng'); // 选择语言
    const [downloadLink, setDownloadLink] = useState(''); // 服务端返回的文件链接
    const [isLoading, setIsLoading] = useState(false); // 是否加载中
    const navigate = useNavigate();

    const handleUpload = async (e) => {
        e.preventDefault();

        if (!file) {
            alert('请选择一个文件');
            return;
        }

        const formData = new FormData();
        formData.append('document', file); // 文件
        formData.append('lang', lang);    // 语言

        setIsLoading(true); // 开始加载

        try {
            const { data } = await uploadPDF(formData);
            setDownloadLink(data.data); // 服务端返回的下载链接
        } catch (error) {
            if (error.response?.status === 401) {
                alert('未授权，请登录');
                navigate('/login');
            } else {
                alert('上传失败，请稍后重试');
            }
        } finally {
            setIsLoading(false); // 停止加载
        }
    };

    return (
        <form onSubmit={handleUpload}>
            <h2>PDF 翻译</h2>
            {/* 自定义文件选择按钮 */}
            <label htmlFor="file-upload" className="file-label">
                选择文件
            </label>
            <div className="spacer"></div>
            <input
                id="file-upload"
                type="file"
                accept=".pdf"
                onChange={(e) => {
                    setFile(e.target.files[0]);
                }}
                disabled={isLoading} // 禁用上传输入框
            />
            {/* 显示选中文件名 */}
            {file && <p className="file-name">已选择：{file.name}</p>}
            <select
                value={lang}
                onChange={(e) => setLang(e.target.value)}
                disabled={isLoading} // 禁用语言选择
            >
                <option value="eng">英文</option>
                <option value="chi_sim">中文</option>
                {/* 添加其他语言选项 */}
            </select>
            <button type="submit" disabled={isLoading}>
                {isLoading ? '处理中...' : '提交'}
            </button>
            {isLoading && <p>处理中，请稍候...</p>} {/* 显示加载提示 */}
            {downloadLink && (
                <div>
                    <p>翻译完成！</p>
                    <a href={downloadLink} target="_blank" rel="noopener noreferrer">下载翻译文件</a>
                </div>
            )}
        </form>
    );
};

export default Translate;
