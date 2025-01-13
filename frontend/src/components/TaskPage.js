import React, {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {fetchTasks} from '../api';

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
                setTasks(Object.entries(response.data.data));
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
        <div className="task-table-container">
            <h2>Your Tasks</h2>
            {tasks.length === 0 ? (
                <p className="no-tasks">No tasks found or loading failed</p>
            ) : (
                <table>
                    <thead>
                    <tr>
                        <th>Task ID</th>
                        <th>Status</th>
                        <th>Download</th>
                    </tr>
                    </thead>
                    <tbody>
                    {tasks.map(([taskId, task]) => (
                        <tr key={taskId}>
                            <td>{taskId}</td>
                            <td>{statusMap[task.status] || 'Unknown Status'}</td>
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
    );
};

export default TaskPage;
