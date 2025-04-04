import React, { useState } from 'react';
import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { Link } from "react-router-dom";
import { FaUser, FaFileAlt, FaKey, FaHistory, FaSearch } from "react-icons/fa";

/**
 * 个人中心页面
 * 显示用户个人信息、数据上传历史及提供授权管理入口
 */
export default function ProfilePage() {
  // 模拟用户数据
  const [userData] = useState({
    username: '张医生',
    role: '医生',
    hospital: '协和医院',
    department: '心内科',
    email: 'zhang@example.com',
    joinDate: '2023-05-15',
  });

  // 模拟用户上传的数据
  const [userUploads] = useState([
    {
      id: 1,
      fileName: '心电图数据.ecg',
      dataType: '影像数据',
      uploadTime: '2024-06-10 14:30',
      status: '已授权',
    },
    {
      id: 2,
      fileName: '患者病历记录.pdf',
      dataType: '电子病历',
      uploadTime: '2024-06-08 09:15',
      status: '未授权',
    },
    {
      id: 3,
      fileName: '基因测序结果.xml',
      dataType: '基因组数据',
      uploadTime: '2024-06-01 16:45',
      status: '已授权',
    },
  ]);

  return (
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 text-white">
      {/* 导航栏 */}
      <nav className="flex justify-between items-center p-6 bg-opacity-80 bg-gray-800 shadow-lg w-full">
        <h1 className="text-2xl font-bold">MedCross</h1>
        <div className="space-x-6">
          <Link to="/" className="hover:text-blue-400">首页</Link>
          <Link to="/data-upload" className="hover:text-blue-400">数据上传</Link>
          <Link to="/data-query" className="hover:text-blue-400">数据查询</Link>
          <Link to="/profile" className="text-blue-400">个人中心</Link>
          <Link to="/blockchain-records" className="hover:text-blue-400">区块链记录</Link>
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        <h2 className="text-3xl font-bold text-center text-blue-400 mb-8">个人中心</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* 左侧个人信息卡片 */}
          <Card className="bg-gray-800 p-6 rounded-xl shadow-lg">
            <div className="flex flex-col items-center">
              <div className="w-24 h-24 rounded-full bg-blue-500 flex items-center justify-center mb-4">
                <FaUser className="text-4xl" />
              </div>
              <h3 className="text-xl font-bold">{userData.username}</h3>
              <p className="text-blue-400">{userData.role}</p>
            </div>
            
            <div className="mt-6 space-y-2">
              <p><span className="text-gray-400">医院：</span>{userData.hospital}</p>
              <p><span className="text-gray-400">科室：</span>{userData.department}</p>
              <p><span className="text-gray-400">邮箱：</span>{userData.email}</p>
              <p><span className="text-gray-400">加入时间：</span>{userData.joinDate}</p>
            </div>
            
            <Button className="w-full mt-6 bg-blue-500 hover:bg-blue-600">
              编辑个人信息
            </Button>
          </Card>
          
          {/* 右侧内容区域 */}
          <div className="md:col-span-2 space-y-8">
            {/* 快捷功能区 */}
            <div className="grid grid-cols-1 sm:grid-cols-4 gap-4">
              <Link to="/data-upload">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer">
                  <div className="flex items-center p-4">
                    <FaFileAlt className="text-blue-400 text-2xl mr-4" />
                    <CardContent>上传新数据</CardContent>
                  </div>
                </Card>
              </Link>
              
              <Link to="/data-query">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer">
                  <div className="flex items-center p-4">
                    <FaSearch className="text-blue-400 text-2xl mr-4" />
                    <CardContent>查询数据</CardContent>
                  </div>
                </Card>
              </Link>
              
              <Link to="/auth-management">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer">
                  <div className="flex items-center p-4">
                    <FaKey className="text-blue-400 text-2xl mr-4" />
                    <CardContent>授权管理</CardContent>
                  </div>
                </Card>
              </Link>
              
              <Link to="/blockchain-records">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer">
                  <div className="flex items-center p-4">
                    <FaHistory className="text-blue-400 text-2xl mr-4" />
                    <CardContent>区块链记录</CardContent>
                  </div>
                </Card>
              </Link>
            </div>
            
            {/* 数据上传历史 */}
            <div>
              <h3 className="text-xl font-bold mb-4">我的数据</h3>
              <div className="bg-gray-800 rounded-xl overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-700">
                    <tr>
                      <th className="py-3 px-4 text-left">文件名</th>
                      <th className="py-3 px-4 text-left">数据类型</th>
                      <th className="py-3 px-4 text-left">上传时间</th>
                      <th className="py-3 px-4 text-left">状态</th>
                      <th className="py-3 px-4 text-left">操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    {userUploads.map((item) => (
                      <tr key={item.id} className="border-t border-gray-700">
                        <td className="py-3 px-4">{item.fileName}</td>
                        <td className="py-3 px-4">{item.dataType}</td>
                        <td className="py-3 px-4">{item.uploadTime}</td>
                        <td className="py-3 px-4">
                          <span className={`px-2 py-1 rounded-full text-xs ${item.status === '已授权' ? 'bg-green-800 text-green-200' : 'bg-yellow-800 text-yellow-200'}`}>
                            {item.status}
                          </span>
                        </td>
                        <td className="py-3 px-4">
                          <Link to={`/data-detail/${item.id}`} className="text-blue-400 hover:underline mr-3">
                            查看
                          </Link>
                          <Link to={`/auth-management/${item.id}`} className="text-blue-400 hover:underline">
                            授权
                          </Link>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}