import React, { useState } from 'react';
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Card } from "../components/ui/card";
import { motion } from "framer-motion";

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
    <div className="min-h-screen w-full bg-gradient-to-b from-blue-900 to-gray-900 flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="w-full max-w-md p-8 bg-gray-800 text-white shadow-xl rounded-xl">
          <h2 className="text-3xl font-bold text-center text-blue-400 mb-6">
            {isLogin ? '登录' : '注册'}
          </h2>
          
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">用户名</label>
              <Input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleInputChange}
                required
                className="w-full bg-gray-700 border-gray-600"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium mb-1">密码</label>
              <Input
                type="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                required
                className="w-full bg-gray-700 border-gray-600"
              />
            </div>
            
            {!isLogin && (
              <>
                <div>
                  <label className="block text-sm font-medium mb-1">确认密码</label>
                  <Input
                    type="password"
                    name="confirmPassword"
                    value={formData.confirmPassword}
                    onChange={handleInputChange}
                    required
                    className="w-full bg-gray-700 border-gray-600"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">电子邮箱</label>
                  <Input
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleInputChange}
                    required
                    className="w-full bg-gray-700 border-gray-600"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium mb-1">角色</label>
                  <select
                    name="role"
                    value={formData.role}
                    onChange={handleInputChange}
                    className="w-full p-2 rounded-md bg-gray-700 border-gray-600"
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
              className="w-full bg-blue-500 hover:bg-blue-600 mt-6"
            >
              {isLogin ? '登录' : '注册'}
            </Button>
          </form>
          
          <div className="text-center mt-4">
            <button
              type="button"
              onClick={() => setIsLogin(!isLogin)}
              className="text-blue-400 hover:underline"
            >
              {isLogin ? '没有账号？点击注册' : '已有账号？点击登录'}
            </button>
          </div>
        </Card>
      </motion.div>
    </div>
  );
}