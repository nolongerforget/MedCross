/**
 * API服务
 * 处理与后端的所有API通信
 */

const API_BASE_URL = 'http://localhost:8000/api';

/**
 * 基础fetch封装，处理错误和响应
 */
async function fetchAPI(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // 默认请求头
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // 如果有token，添加到请求头
  const token = localStorage.getItem('authToken');
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch(url, {
      ...options,
      headers,
    });

    // 检查响应状态
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `请求失败: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('API请求错误:', error);
    throw error;
  }
}

/**
 * 认证相关API
 */
export const authAPI = {
  // 用户登录
  login: async (username, password) => {
    return fetchAPI('/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
  },

  // 用户注册
  register: async (userData) => {
    return fetchAPI('/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  },

  // 获取当前用户信息
  getCurrentUser: async () => {
    return fetchAPI('/user');
  },
};

/**
 * 数据相关API
 */
export const dataAPI = {
  // 查询数据
  queryData: async (params) => {
    const queryString = new URLSearchParams(params).toString();
    return fetchAPI(`/query?${queryString}`);
  },

  // 上传数据
  uploadData: async (data) => {
    return fetchAPI('/upload', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // 获取数据类型列表
  getDataTypes: async () => {
    return fetchAPI('/data-types');
  },

  // 获取数据详情
  getDataDetail: async (id) => {
    return fetchAPI(`/data/${id}`);
  },

  // 获取统计数据
  getStatistics: async () => {
    // 模拟统计数据，实际项目中应从后端获取
    return {
      totalRecords: 1256,
      ethereumRecords: 723,
      fabricRecords: 533,
      dataTypes: {
        '影像数据': 482,
        '电子病历': 345,
        '基因组数据': 128,
        '处方数据': 201,
        '检验报告': 100,
      },
      monthlyGrowth: [
        { month: '1月', count: 45 },
        { month: '2月', count: 52 },
        { month: '3月', count: 78 },
        { month: '4月', count: 110 },
        { month: '5月', count: 145 },
        { month: '6月', count: 189 },
      ]
    };
  }
};