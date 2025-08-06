const API_BASE = '/api';

// 显示加载动画
function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

// 隐藏加载动画
function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

// 显示响应结果
function showResponse(elementId, data, isSuccess = true) {
    const element = document.getElementById(elementId);
    element.style.display = 'block';
    element.className = `response ${isSuccess ? 'success' : 'error'}`;
    element.textContent = JSON.stringify(data, null, 2);
}

// 获取数据库配置
function getDatabaseConfig() {
    return {
        host: document.getElementById('dbHost').value,
        port: parseInt(document.getElementById('dbPort').value),
        user: document.getElementById('dbUser').value,
        password: document.getElementById('dbPassword').value,
        dbname: document.getElementById('dbName').value
    };
}

// 检查数据库状态
async function checkDatabaseStatus() {
    try {
        showLoading();
        const response = await fetch(`${API_BASE}/database/status`);
        const data = await response.json();
        
        const statusElement = document.getElementById('dbStatus');
        if (data.success && data.data.connected) {
            statusElement.innerHTML = `
                <span class="status-indicator status-connected"></span>
                已连接 - 活跃连接: ${data.data.open_connections}
            `;
        } else {
            statusElement.innerHTML = `
                <span class="status-indicator status-disconnected"></span>
                ${data.data ? data.data.status : '未连接'}
            `;
        }
    } catch (error) {
        document.getElementById('dbStatus').innerHTML = `
            <span class="status-indicator status-disconnected"></span>
            连接错误: ${error.message}
        `;
    } finally {
        hideLoading();
    }
}

// 创建数据库
async function createDatabase() {
    try {
        showLoading();
        const config = getDatabaseConfig();
        
        const response = await fetch(`${API_BASE}/database/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
        });
        
        const data = await response.json();
        showResponse('dbResponse', data, data.success);
        
        if (data.success) {
            setTimeout(checkDatabaseStatus, 1000);
        }
    } catch (error) {
        showResponse('dbResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}

// 连接数据库
async function connectDatabase() {
    try {
        showLoading();
        const config = getDatabaseConfig();
        
        const response = await fetch(`${API_BASE}/database/connect`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
        });
        
        const data = await response.json();
        showResponse('dbResponse', data, data.success);
        
        setTimeout(checkDatabaseStatus, 1000);
    } catch (error) {
        showResponse('dbResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}

// 创建用户
async function createUser() {
    try {
        showLoading();
        const userData = {
            name: document.getElementById('userName').value,
            email: document.getElementById('userEmail').value,
            age: parseInt(document.getElementById('userAge').value) || 0
        };

        if (!userData.name || !userData.email) {
            throw new Error('请填写姓名和邮箱');
        }
        
        const response = await fetch(`${API_BASE}/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        
        const data = await response.json();
        showResponse('userResponse', data, data.success);
        
        if (data.success) {
            clearUserForm();
            setTimeout(getAllUsers, 1000);
        }
    } catch (error) {
        showResponse('userResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}

// 发表留言
async function addContent() {
    showLoading();
    const ctt = document.getElementById('userContent').value;
    if (!ctt) {
        hideLoading();
        showResponse('contentResponse', { error: '留言内容不能为空' }, false);
        return;
    } else {
        hideLoading();
        const response = await fetch(`${API_BASE}/addContent`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ user_id: 'todo', content: ctt })
        })
        showResponse('contentResponse', {success: `留言内容: ${ctt}` }, true);
        document.getElementById('userContent').value = ''; // 清空输入框
        setTimeout(getAllContent, 1000); // 刷新留言列表
        return
    }
}

// 获取所有用户
async function getAllUsers() {
    try {
        showLoading();
        const response = await fetch(`${API_BASE}/users`);
        const data = await response.json();
        
        //showResponse('userResponse', data, data.success);
        
        if (data.success && Array.isArray(data.data)) {
            displayUsers(data.data);
        }
    } catch (error) {
        showResponse('userResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}
async function getAllContent() {
    try {
        showLoading();
        const response = await fetch(`${API_BASE}/getAllContent`);
        const data = await response.json();

        if (data.success && Array.isArray(data.data)) {
            displayContent(data.data);
        }
    } catch (error) {
        showResponse('userContentInfo', { error: error.message }, false);
    } finally {
        hideLoading();
    }

}

// 显示用户列表
function displayUsers(users) {
    const container = document.getElementById('usersTableContainer');
    const tbody = document.getElementById('usersTableBody');
    
    tbody.innerHTML = '';
    
    users.forEach(user => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${user.id}</td>
            <td>${user.name}</td>
            <td>${user.email}</td>
            <td>${user.age}</td>
            <td>${new Date(user.created_at).toLocaleString()}</td>
            <td>
                <button onclick="editUser(${user.id})" class="btn-success" style="margin: 2px; padding: 8px 12px; font-size: 14px;">编辑</button>
                <button onclick="deleteUser(${user.id})" class="btn-danger" style="margin: 2px; padding: 8px 12px; font-size: 14px;">删除</button>
            </td>
        `;
        tbody.appendChild(row);
    });
    
    container.style.display = users.length > 0 ? 'block' : 'none';
}

function displayContent (cttList) {
    const container = document.getElementById('contentTableContainer');
    container.style.display = cttList.length > 0 ? 'block' : 'none';
    const tbody = document.getElementById('contentTableBody');

    tbody.innerHTML = '';
    cttList.forEach(ctt => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${ctt.user_ip}</td>
            <td>${ctt.content}</td>
            <td>${ctt.created_at}</td>
            <td>
                <button onclick="deleteContent(${ctt.msg_id})" class="btn-danger" style="margin: 2px; padding: 8px 12px; font-size: 14px;">删除</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// 删除用户
async function deleteUser(id) {
    if (!confirm('确定要删除这个用户吗？')) {
        return;
    }

    try {
        showLoading();
        const response = await fetch(`${API_BASE}/users/${id}`, {
            method: 'DELETE'
        });
        
        const data = await response.json();
        showResponse('userResponse', data, data.success);
        
        if (data.success) {
            setTimeout(getAllUsers, 1000);
        }
    } catch (error) {
        showResponse('userResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}

async function deleteContent(msg_id) {
    if (!confirm('确定要删除这条留言吗？')) {
        return;
    }
    try {
        showLoading();
        const resp = await fetch(`${API_BASE}/removeContent/${msg_id}`, {
            method: 'DELETE'
        });

        const data = await resp.json();
        showResponse('contentResponse', data, data.success);
        if (data.success) {
            setTimeout(getAllContent, 1000);
        }
    } catch (error) {
        showResponse('contentResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}

// 编辑用户
async function editUser(id) {
    const name = prompt('请输入新的姓名:');
    const email = prompt('请输入新的邮箱:');
    const age = prompt('请输入新的年龄:');
    
    if (!name && !email && !age) {
        return;
    }

    try {
        showLoading();
        const updateData = {};
        if (name) updateData.name = name;
        if (email) updateData.email = email;
        if (age) updateData.age = parseInt(age);
        
        const response = await fetch(`${API_BASE}/users/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updateData)
        });
        
        const data = await response.json();
        showResponse('userResponse', data, data.success);
        
        if (data.success) {
            setTimeout(getAllUsers, 1000);
        }
    } catch (error) {
        showResponse('userResponse', { error: error.message }, false);
    } finally {
        hideLoading();
    }
}

// 清空用户表单
function clearUserForm() {
    document.getElementById('userName').value = '';
    document.getElementById('userEmail').value = '';
    document.getElementById('userAge').value = '';
}

// 页面加载完成后检查数据库状态
document.addEventListener('DOMContentLoaded', function() {
    checkDatabaseStatus();
});