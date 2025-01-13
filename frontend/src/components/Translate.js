import React, {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {fetchUserInfo, uploadPDF} from '../api';
import {logout} from "../utils";

const Translate = () => {
    const [file, setFile] = useState(null);
    const [lang, setLang] = useState('eng');
    const [isLoading, setIsLoading] = useState(false);
    const [isSidebarVisible, setIsSidebarVisible] = useState(false);
    const [userInfo, setUserInfo] = useState({username: '', balance: 0});
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
            const response = await uploadPDF(formData);
            if (response.status === 200) {
                alert('Submit successfully. You can view the task status on the task page.');
            } else if (response.data && response.data.error) {
                alert(`Error: ${response.data.error}`);
            } else {
                alert('Unexpected response from server.');
            }
        } catch (error) {
            alert('Upload failed: Unable to connect to the server.');
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
                console.error('Failed to fetch user info:', error);
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
            </form>

            <div style={{marginTop: '10px', display: 'flex', alignItems: 'center', gap: '10px'}}>
                <button onClick={handleLogout}>Logout</button>
                <button onClick={() => navigate('/tasks')}>View Tasks</button>
                <button onClick={toggleSidebar} className="user-info-button">User Info</button>
            </div>

            {isSidebarVisible && (
                <div className="sidebar">
                    <h3>User Info</h3>
                    <p>Username: {userInfo.username}</p>
                    <p>Balance: {userInfo.balance}</p>
                </div>
            )}
        </div>
    );
};

export default Translate;
