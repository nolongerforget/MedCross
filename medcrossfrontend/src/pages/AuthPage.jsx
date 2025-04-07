import React, { useState } from 'react';
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card } from "../components/ui/card";
import { motion } from "framer-motion";
import { Link } from "react-router-dom";

/**
 * 用户认证页面 - 包含登录和注册功能
 * 用于医疗数据跨链共享平台的用户身份验证
 */
export default function AuthPage() {
  // 状态管理
  const [isLogin, setIsLogin] = useState(true); // 控制显示登录还是注册表单
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
    role: '医生' // 默认角色
  });

  // 处理表单输入变化
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  // 处理表单提交
  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (isLogin) {
      // 登录逻辑
      console.log('登录信息:', {
        username: formData.username,
        password: formData.password
      });
      // 这里应该调用API进行身份验证
    } else {
      // 注册逻辑
      if (formData.password !== formData.confirmPassword) {
        alert('两次输入的密码不一致！');
        return;
      }
      console.log('注册信息:', formData);
      // 这里应该调用API进行注册
    }
  };

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
          <Link to="/auth" className="bg-blue-500 hover:bg-blue-600 px-3 py-1 rounded-md">登录/注册</Link>
        </div>
      </nav>
      
      <div className="flex items-center justify-center p-8 min-h-[calc(100vh-80px)]">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="w-full max-w-xl"
        >
          <Card className="p-8 bg-gray-800 text-white shadow-xl rounded-xl border border-gray-700">
            <h2 className="text-3xl font-bold text-center text-blue-400 mb-8">
              {isLogin ? '登录账号' : '注册新账号'}
            </h2>
            
            <form onSubmit={handleSubmit} className="space-y-5">
              <div>
                <label className="block text-sm font-medium mb-2">用户名</label>
                <Input
                  type="text"
                  name="username"
                  value={formData.username}
                  onChange={handleInputChange}
                  required
                  placeholder="请输入用户名"
                  className="bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2">密码</label>
                <Input
                  type="password"
                  name="password"
                  value={formData.password}
                  onChange={handleInputChange}
                  required
                  placeholder="请输入密码"
                  className="bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                />
              </div>
              
              {!isLogin && (
                <>
                  <div>
                    <label className="block text-sm font-medium mb-2">确认密码</label>
                    <Input
                      type="password"
                      name="confirmPassword"
                      value={formData.confirmPassword}
                      onChange={handleInputChange}
                      required
                      placeholder="请再次输入密码"
                      className="bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                    />
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium mb-2">电子邮箱</label>
                    <Input
                      type="email"
                      name="email"
                      value={formData.email}
                      onChange={handleInputChange}
                      required
                      placeholder="请输入电子邮箱"
                      className="bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                    />
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium mb-2">角色</label>
                    <select
                      name="role"
                      value={formData.role}
                      onChange={handleInputChange}
                      className="w-full p-3 rounded-md bg-gray-700 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="医生">医生</option>
                      <option value="研究人员">研究人员</option>
                      <option value="医院管理员">医院管理员</option>
                      <option value="患者">患者</option>
                    </select>
                  </div>
                </>
              )}
              
              <Button 
                type="submit" 
                className="w-full bg-blue-500 hover:bg-blue-600 mt-8 py-3 text-lg font-medium"
              >
                {isLogin ? '登录' : '注册'}
              </Button>
            </form>
            
            <div className="text-center mt-6">
              <button
                type="button"
                onClick={() => setIsLogin(!isLogin)}
                className="text-blue-400 hover:underline text-base"
              >
                {isLogin ? '没有账号？点击注册' : '已有账号？点击登录'}
              </button>
            </div>
          </Card>
        </motion.div>
      </div>
    </div>
  );
}