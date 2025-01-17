import React, {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {fetchTasks} from '../api';
import {Tooltip} from "react-tooltip";
import {formatTimestamp} from "../utils"

const TaskPage = () => {
    const [tasks, setTasks] = useState([]);
    const navigate = useNavigate();

    const statusMap = {
        0: 'Task Received',
        1: 'Translating',
        2: 'Uploading',
        3: 'Done',
        9: 'Error',
    };

    useEffect(() => {
        let intervalId;

        const loadTasks = async () => {
            try {
                const response = await fetchTasks();
                // 获取 tasks 数据并直接存储为值数组
                setTasks(Object.values(response.data.data));
            } catch (error) {
                alert('Failed to load tasks');
            }
        };

        // 初始化加载任务
        loadTasks();

        // 定时刷新任务状态
        intervalId = setInterval(() => {
            loadTasks();
        }, 5000); // 每5秒刷新一次

        return () => {
            clearInterval(intervalId); // 清除定时器
        };
    }, []);

    return (
        <>
        <div className="task-table-container">
            <h2>Your Tasks</h2>
            {tasks.length === 0 ? (
                <p className="no-tasks">No tasks found or loading failed</p>
            ) : (
                <table>
                    <thead>
                    <tr>
                        <th>Filename</th>
                        <th>Status</th>
                        <th>Upload Time</th>
                        <th>Download</th>
                    </tr>
                    </thead>
                    <tbody>
                    {tasks.map((task, index) => (
                        <tr key={index}>
                            <td>
                                {task.filename.length > 10
                                    ? `${task.filename.substring(0, 10)}...`
                                    : task.filename}
                            </td>
                            <td>{statusMap[task.status] || 'Unknown Status'}</td>
                            <td>{formatTimestamp(task.created_at) || 'Unknown Created Time'}</td>
                            <td>
                                {task.link ? (
                                    <a
                                        href={task.link}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="download-link"
                                    >
                                        Download
                                    </a>
                                ) : (
                                    'N/A'
                                )}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
            <div style={{marginTop: '20px'}}>
                <button onClick={() => navigate('/')} className="back-button">
                    Back to Translate
                </button>
            </div>

        </div>
            <div className="more-info-container">
                <span
                    data-tooltip-id="info-tooltip"
                    className="more-info"
                >
                    More Info
                </span>
                <Tooltip id="info-tooltip" place="top" type="dark" effect="solid">
                    This page is automatically refreshed.
                </Tooltip>
        </div>
        </>
    )
        ;
};

export default TaskPage;
