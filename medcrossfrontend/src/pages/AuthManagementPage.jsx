import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card } from "../components/ui/card";
import { FaKey, FaSearch, FaPlus, FaTimes } from "react-icons/fa";

/**
 * 数据授权管理页面
 * 用于管理医疗数据的授权，包括查看现有授权、添加新授权和撤销授权
 */
export default function AuthManagementPage() {
  const { id } = useParams(); // 可选参数，如果有ID则是特定数据的授权管理
  const [searchTerm, setSearchTerm] = useState('');
  const [showAddModal, setShowAddModal] = useState(false);
  const [selectedData, setSelectedData] = useState(null);
  const [loading, setLoading] = useState(true);
  
  // 模拟数据列表
  const [dataList, setDataList] = useState([]);
  
  // 模拟授权接收方列表
  const [recipients] = useState([
    { id: 1, name: '李医生', hospital: '人民医院', department: '神经内科' },
    { id: 2, name: '王研究员', hospital: '医学研究中心', department: '基因研究室' },
    { id: 3, name: '赵医生', hospital: '第三医院', department: '心胸外科' },
    { id: 4, name: '医学数据分析中心', hospital: '医科大学', department: '数据科学部' },
  ]);

  // 模拟从API获取数据列表
  useEffect(() => {
    setTimeout(() => {
      const data = [
        {
          id: '1',
          fileName: '心电图数据.ecg',
          dataType: '影像数据',
          uploadTime: '2024-06-10 14:30',
          authorizations: [
            { id: 1, recipient: '李医生', hospital: '人民医院', time: '2024-06-12 09:30', status: '已授权' },
          ]
        },
        {
          id: '2',
          fileName: '患者病历记录.pdf',
          dataType: '电子病历',
          uploadTime: '2024-06-08 09:15',
          authorizations: []
        },
        {
          id: '3',
          fileName: '基因测序结果.xml',
          dataType: '基因组数据',
          uploadTime: '2024-06-01 16:45',
          authorizations: [
            { id: 2, recipient: '医学研究中心', hospital: '医科大学', time: '2024-06-15 14:20', status: '已授权' },
          ]
        },
      ];
      
      // 如果URL中有ID参数，则只显示该ID对应的数据
      if (id) {
        const filtered = data.filter(item => item.id === id);
        setDataList(filtered);
        if (filtered.length > 0) {
          setSelectedData(filtered[0]);
        }
      } else {
        setDataList(data);
      }
      
      setLoading(false);
    }, 1000); // 模拟加载延迟
  }, [id]);

  // 处理搜索
  const handleSearch = () => {
    // 实际应用中这里应该调用API进行搜索
    console.log('搜索:', searchTerm);
  };

  // 处理添加授权
  const handleAddAuthorization = (dataId, recipientId) => {
    console.log('添加授权:', dataId, recipientId);
    // 实际应用中这里应该调用API添加授权
    
    // 模拟添加授权成功
    setDataList(prevList => {
      return prevList.map(data => {
        if (data.id === dataId) {
          const recipient = recipients.find(r => r.id === recipientId);
          const newAuth = {
            id: Date.now(), // 使用时间戳作为临时ID
            recipient: recipient.name,
            hospital: recipient.hospital,
            time: new Date().toLocaleString(),
            status: '已授权'
          };
          return {
            ...data,
            authorizations: [...data.authorizations, newAuth]
          };
        }
        return data;
      });
    });
    
    setShowAddModal(false);
  };

  // 处理撤销授权
  const handleRevokeAuthorization = (dataId, authId) => {
    console.log('撤销授权:', dataId, authId);
    // 实际应用中这里应该调用API撤销授权
    
    // 模拟撤销授权成功
    setDataList(prevList => {
      return prevList.map(data => {
        if (data.id === dataId) {
          return {
            ...data,
            authorizations: data.authorizations.filter(auth => auth.id !== authId)
          };
        }
        return data;
      });
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white flex items-center justify-center">
        <p className="text-xl">加载中...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white">
      {/* 导航栏 */}
      <nav className="flex justify-between items-center p-6 bg-opacity-80 bg-gray-800 shadow-lg w-full">
        <h1 className="text-2xl font-bold">MedCross</h1>
        <div className="space-x-6">
          <Link to="/" className="hover:text-blue-400">首页</Link>
          <Link to="/data-upload" className="hover:text-blue-400">数据上传</Link>
          <Link to="/data-query" className="hover:text-blue-400">数据查询</Link>
          <Link to="/profile" className="hover:text-blue-400">个人中心</Link>
          <Link to="/blockchain-records" className="hover:text-blue-400">区块链记录</Link>
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        {/* 返回按钮 */}
        <div className="mb-6">
          <Link to="/profile" className="text-blue-400 hover:underline flex items-center">
            <span>← 返回个人中心</span>
          </Link>
        </div>

        <h2 className="text-3xl font-bold text-blue-400 mb-6">
          {id ? '数据授权管理' : '我的数据授权管理'}
        </h2>

        {/* 搜索栏 - 仅在查看所有数据时显示 */}
        {!id && (
          <div className="mb-6 flex">
            <Input
              type="text"
              placeholder="搜索数据文件名或类型"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="flex-grow bg-gray-700 border-gray-600"
            />
            <Button onClick={handleSearch} className="ml-2 bg-blue-500 hover:bg-blue-600">
              <FaSearch className="mr-2" /> 搜索
            </Button>
          </div>
        )}

        {/* 数据列表 */}
        {dataList.length > 0 ? (
          dataList.map((data) => (
            <Card key={data.id} className="mb-6 bg-gray-800 p-6">
              <div className="flex justify-between items-center mb-4">
                <div>
                  <h3 className="text-xl font-bold">{data.fileName}</h3>
                  <p className="text-gray-400">{data.dataType} • 上传于 {data.uploadTime}</p>
                </div>
                <Button 
                  onClick={() => {
                    setSelectedData(data);
                    setShowAddModal(true);
                  }} 
                  className="bg-green-600 hover:bg-green-700"
                >
                  <FaPlus className="mr-2" /> 添加授权
                </Button>
              </div>

              {/* 授权列表 */}
              <div className="mt-4">
                <h4 className="font-bold mb-2">授权列表</h4>
                {data.authorizations.length > 0 ? (
                  <div className="bg-gray-700 rounded-lg overflow-hidden">
                    <table className="w-full">
                      <thead className="bg-gray-800">
                        <tr>
                          <th className="py-3 px-4 text-left">接收方</th>
                          <th className="py-3 px-4 text-left">所属机构</th>
                          <th className="py-3 px-4 text-left">授权时间</th>
                          <th className="py-3 px-4 text-left">状态</th>
                          <th className="py-3 px-4 text-left">操作</th>
                        </tr>
                      </thead>
                      <tbody>
                        {data.authorizations.map((auth) => (
                          <tr key={auth.id} className="border-t border-gray-600">
                            <td className="py-3 px-4">{auth.recipient}</td>
                            <td className="py-3 px-4">{auth.hospital}</td>
                            <td className="py-3 px-4">{auth.time}</td>
                            <td className="py-3 px-4">
                              <span className="px-2 py-1 rounded-full text-xs bg-green-800 text-green-200">
                                {auth.status}
                              </span>
                            </td>
                            <td className="py-3 px-4">
                              <button 
                                onClick={() => handleRevokeAuthorization(data.id, auth.id)}
                                className="text-red-400 hover:underline flex items-center"
                              >
                                <FaTimes className="mr-1" /> 撤销
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <p className="text-gray-400">暂无授权记录</p>
                )}
              </div>
            </Card>
          ))
        ) : (
          <div className="text-center py-10">
            <FaKey className="text-5xl text-gray-600 mx-auto mb-4" />
            <p className="text-xl text-gray-400">暂无数据授权记录</p>
          </div>
        )}
      </div>

      {/* 添加授权模态框 */}
      {showAddModal && selectedData && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-gray-800 rounded-xl p-6 w-full max-w-md">
            <h3 className="text-xl font-bold mb-4">添加授权</h3>
            <p className="mb-4">为 <span className="font-bold">{selectedData.fileName}</span> 添加新的授权</p>
            
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">选择接收方</label>
              <div className="space-y-2 max-h-60 overflow-y-auto">
                {recipients.map((recipient) => (
                  <div 
                    key={recipient.id} 
                    className="bg-gray-700 p-3 rounded-lg cursor-pointer hover:bg-gray-600"
                    onClick={() => handleAddAuthorization(selectedData.id, recipient.id)}
                  >
                    <div className="font-bold">{recipient.name}</div>
                    <div className="text-sm text-gray-400">{recipient.hospital} • {recipient.department}</div>
                  </div>
                ))}
              </div>
            </div>
            
            <div className="flex justify-end mt-6">
              <Button 
                onClick={() => setShowAddModal(false)} 
                className="bg-gray-600 hover:bg-gray-700 mr-2"
              >
                取消
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}