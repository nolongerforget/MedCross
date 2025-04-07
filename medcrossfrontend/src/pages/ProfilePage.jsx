import React, { useState } from 'react';
import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { Link } from "react-router-dom";
import { FaUser, FaFileAlt, FaKey, FaSearch, FaEdit, FaShieldAlt } from "react-icons/fa";

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
        </div>
      </nav>

      <div className="container mx-auto py-8 px-4">
        <h2 className="text-3xl font-bold text-center text-blue-400 mb-8">个人中心</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* 左侧个人信息卡片 */}
          <Card className="bg-gray-800 p-6 rounded-xl shadow-lg border border-gray-700">
            <div className="flex flex-col items-center">
              <div className="w-28 h-28 rounded-full bg-gradient-to-r from-blue-500 to-blue-700 flex items-center justify-center mb-4 shadow-lg">
                <FaUser className="text-4xl" />
              </div>
              <h3 className="text-2xl font-bold">{userData.username}</h3>
              <p className="text-blue-400 text-lg">{userData.role}</p>
            </div>
            
            <div className="mt-8 space-y-3">
              <div className="flex items-center">
                <span className="text-gray-400 w-24">医院：</span>
                <span className="text-white">{userData.hospital}</span>
              </div>
              <div className="flex items-center">
                <span className="text-gray-400 w-24">科室：</span>
                <span className="text-white">{userData.department}</span>
              </div>
              <div className="flex items-center">
                <span className="text-gray-400 w-24">邮箱：</span>
                <span className="text-white">{userData.email}</span>
              </div>
              <div className="flex items-center">
                <span className="text-gray-400 w-24">加入时间：</span>
                <span className="text-white">{userData.joinDate}</span>
              </div>
            </div>
            
            <Button className="w-full mt-8 bg-blue-500 hover:bg-blue-600 py-2.5 flex items-center justify-center gap-2">
              <FaEdit className="text-sm" /> 编辑个人信息
            </Button>
          </Card>
          
          {/* 右侧内容区域 */}
          <div className="md:col-span-2 space-y-8">
            {/* 快捷功能区 */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <Link to="/data-upload">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer border border-gray-700 hover:border-blue-500">
                  <div className="flex items-center p-5">
                    <div className="bg-blue-500 bg-opacity-20 p-3 rounded-full mr-4">
                      <FaFileAlt className="text-blue-400 text-2xl" />
                    </div>
                    <CardContent className="text-lg">上传新数据</CardContent>
                  </div>
                </Card>
              </Link>
              
              <Link to="/data-query">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer border border-gray-700 hover:border-blue-500">
                  <div className="flex items-center p-5">
                    <div className="bg-blue-500 bg-opacity-20 p-3 rounded-full mr-4">
                      <FaSearch className="text-blue-400 text-2xl" />
                    </div>
                    <CardContent className="text-lg">查询数据</CardContent>
                  </div>
                </Card>
              </Link>
              
              <Link to="/auth-management">
                <Card className="bg-gray-800 hover:bg-gray-700 transition-colors cursor-pointer border border-gray-700 hover:border-blue-500">
                  <div className="flex items-center p-5">
                    <div className="bg-blue-500 bg-opacity-20 p-3 rounded-full mr-4">
                      <FaKey className="text-blue-400 text-2xl" />
                    </div>
                    <CardContent className="text-lg">授权管理</CardContent>
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