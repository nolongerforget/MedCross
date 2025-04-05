// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title MedicalData
 * @dev 医疗数据管理智能合约，用于存储医疗数据的元数据和管理访问权限
 */
contract MedicalData {
    // 数据结构定义
    struct DataInfo {
        string dataId;          // 数据唯一标识符
        string dataHash;        // 数据哈希值（IPFS或其他存储系统的引用）
        string dataType;        // 数据类型（如：电子病历、影像数据、基因组数据等）
        string description;     // 数据描述
        string[] tags;          // 数据标签
        address owner;          // 数据所有者
        uint256 timestamp;      // 上传时间戳
        bool isConfidential;    // 是否为机密数据
        string fabricTxId;      // Fabric链上对应的交易ID（用于跨链引用）
    }
    
    // 授权结构定义
    struct Authorization {
        string dataId;          // 被授权的数据ID
        address authorizedUser; // 被授权的用户地址
        uint256 startTime;      // 授权开始时间
        uint256 endTime;        // 授权结束时间（0表示永久授权）
        bool isActive;          // 授权是否有效
    }
    
    // 存储所有医疗数据
    mapping(string => DataInfo) private dataRegistry;
    
    // 存储用户拥有的数据ID列表
    mapping(address => string[]) private userDataIds;
    
    // 存储数据的授权信息
    mapping(string => Authorization[]) private dataAuthorizations;
    
    // 事件定义
    event DataUploaded(string dataId, address indexed owner, string dataType, uint256 timestamp);
    event AuthorizationGranted(string dataId, address indexed owner, address indexed authorizedUser, uint256 startTime, uint256 endTime);
    event AuthorizationRevoked(string dataId, address indexed owner, address indexed authorizedUser);
    event DataAccessed(string dataId, address indexed accessor, uint256 timestamp);
    
    /**
     * @dev 上传新的医疗数据
     * @param _dataId 数据ID
     * @param _dataHash 数据哈希
     * @param _dataType 数据类型
     * @param _description 数据描述
     * @param _tags 数据标签
     * @param _isConfidential 是否机密
     * @param _fabricTxId Fabric链上对应的交易ID
     */
    function uploadData(
        string memory _dataId,
        string memory _dataHash,
        string memory _dataType,
        string memory _description,
        string[] memory _tags,
        bool _isConfidential,
        string memory _fabricTxId
    ) public {
        // 确保数据ID不重复
        require(bytes(dataRegistry[_dataId].dataId).length == 0, "Data ID already exists");
        
        // 创建新的数据记录
        DataInfo memory newData = DataInfo({
            dataId: _dataId,
            dataHash: _dataHash,
            dataType: _dataType,
            description: _description,
            tags: _tags,
            owner: msg.sender,
            timestamp: block.timestamp,
            isConfidential: _isConfidential,
            fabricTxId: _fabricTxId
        });
        
        // 存储数据记录
        dataRegistry[_dataId] = newData;
        userDataIds[msg.sender].push(_dataId);
        
        // 触发事件
        emit DataUploaded(_dataId, msg.sender, _dataType, block.timestamp);
    }
    
    /**
     * @dev 授权其他用户访问数据
     * @param _dataId 数据ID
     * @param _authorizedUser 被授权用户地址
     * @param _startTime 授权开始时间
     * @param _endTime 授权结束时间
     */
    function grantAccess(
        string memory _dataId,
        address _authorizedUser,
        uint256 _startTime,
        uint256 _endTime
    ) public {
        // 确保数据存在且调用者是数据所有者
        require(bytes(dataRegistry[_dataId].dataId).length > 0, "Data does not exist");
        require(dataRegistry[_dataId].owner == msg.sender, "Only data owner can grant access");
        
        // 创建新的授权记录
        Authorization memory newAuth = Authorization({
            dataId: _dataId,
            authorizedUser: _authorizedUser,
            startTime: _startTime > 0 ? _startTime : block.timestamp,
            endTime: _endTime,
            isActive: true
        });
        
        // 存储授权记录
        dataAuthorizations[_dataId].push(newAuth);
        
        // 触发事件
        emit AuthorizationGranted(_dataId, msg.sender, _authorizedUser, newAuth.startTime, _endTime);
    }
    
    /**
     * @dev 撤销用户的数据访问权限
     * @param _dataId 数据ID
     * @param _authorizedUser 被授权用户地址
     */
    function revokeAccess(string memory _dataId, address _authorizedUser) public {
        // 确保数据存在且调用者是数据所有者
        require(bytes(dataRegistry[_dataId].dataId).length > 0, "Data does not exist");
        require(dataRegistry[_dataId].owner == msg.sender, "Only data owner can revoke access");
        
        // 查找并撤销授权
        Authorization[] storage auths = dataAuthorizations[_dataId];
        for (uint i = 0; i < auths.length; i++) {
            if (auths[i].authorizedUser == _authorizedUser && auths[i].isActive) {
                auths[i].isActive = false;
                emit AuthorizationRevoked(_dataId, msg.sender, _authorizedUser);
                break;
            }
        }
    }
    
    /**
     * @dev 检查用户是否有权访问数据
     * @param _dataId 数据ID
     * @param _user 用户地址
     * @return 是否有权访问
     */
    function checkAccess(string memory _dataId, address _user) public view returns (bool) {
        // 数据所有者始终有访问权限
        if (dataRegistry[_dataId].owner == _user) {
            return true;
        }
        
        // 检查授权记录
        Authorization[] storage auths = dataAuthorizations[_dataId];
        for (uint i = 0; i < auths.length; i++) {
            if (auths[i].authorizedUser == _user && auths[i].isActive) {
                // 检查授权是否在有效期内
                if (auths[i].startTime <= block.timestamp && 
                    (auths[i].endTime == 0 || auths[i].endTime >= block.timestamp)) {
                    return true;
                }
            }
        }
        
        return false;
    }
    
    /**
     * @dev 记录数据访问事件
     * @param _dataId 数据ID
     */
    function logAccess(string memory _dataId) public {
        // 确保数据存在且调用者有权访问
        require(bytes(dataRegistry[_dataId].dataId).length > 0, "Data does not exist");
        require(checkAccess(_dataId, msg.sender), "Access denied");
        
        // 记录访问事件
        emit DataAccessed(_dataId, msg.sender, block.timestamp);
    }
    
    /**
     * @dev 获取数据详情
     * @param _dataId 数据ID
     * @return 数据详情
     */
    function getDataInfo(string memory _dataId) public view returns (DataInfo memory) {
        // 确保数据存在且调用者有权访问
        require(bytes(dataRegistry[_dataId].dataId).length > 0, "Data does not exist");
        require(checkAccess(_dataId, msg.sender), "Access denied");
        
        return dataRegistry[_dataId];
    }
    
    /**
     * @dev 获取用户拥有的所有数据ID
     * @param _user 用户地址
     * @return 数据ID数组
     */
    function getUserDataIds(address _user) public view returns (string[] memory) {
        return userDataIds[_user];
    }
    
    /**
     * @dev 获取数据的所有授权信息
     * @param _dataId 数据ID
     * @return 授权信息数组
     */
    function getDataAuthorizations(string memory _dataId) public view returns (Authorization[] memory) {
        // 确保数据存在且调用者是数据所有者
        require(bytes(dataRegistry[_dataId].dataId).length > 0, "Data does not exist");
        require(dataRegistry[_dataId].owner == msg.sender, "Only data owner can view authorizations");
        
        return dataAuthorizations[_dataId];
    }
    
    /**
     * @dev 更新数据的Fabric交易ID（用于跨链引用更新）
     * @param _dataId 数据ID
     * @param _fabricTxId 新的Fabric交易ID
     */
    function updateFabricTxId(string memory _dataId, string memory _fabricTxId) public {
        // 确保数据存在且调用者是数据所有者
        require(bytes(dataRegistry[_dataId].dataId).length > 0, "Data does not exist");
        require(dataRegistry[_dataId].owner == msg.sender, "Only data owner can update Fabric TxID");
        
        dataRegistry[_dataId].fabricTxId = _fabricTxId;
    }
}