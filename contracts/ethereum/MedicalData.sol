// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title MedicalData
 * @dev 用于医疗数据上传和查询的智能合约
 */
contract MedicalData {
    // 数据结构定义
    struct Data {
        uint256 id;
        address owner;
        string dataHash;      // IPFS或其他存储系统的哈希值
        string dataType;      // 数据类型（如：影像数据、电子病历等）
        string metadata;      // JSON格式的元数据
        uint256 timestamp;    // 上传时间戳
        string keywords;      // 关键词，用于搜索，以逗号分隔
    }
    
    // 存储所有医疗数据
    Data[] private allData;
    
    // 用户拥有的数据映射
    mapping(address => uint256[]) private userDataIds;
    
    // 数据类型到数据ID的映射
    mapping(string => uint256[]) private typeToDataIds;
    
    // 事件定义
    event DataUploaded(uint256 indexed id, address indexed owner, string dataType, uint256 timestamp);
    
    /**
     * @dev 上传新的医疗数据
     * @param dataHash 数据的哈希值
     * @param dataType 数据类型
     * @param metadata 元数据（JSON格式）
     * @param keywords 关键词，用于搜索
     * @return 新数据的ID
     */
    function uploadData(
        string memory dataHash,
        string memory dataType,
        string memory metadata,
        string memory keywords
    ) public returns (uint256) {
        uint256 newId = allData.length;
        
        Data memory newData = Data({
            id: newId,
            owner: msg.sender,
            dataHash: dataHash,
            dataType: dataType,
            metadata: metadata,
            timestamp: block.timestamp,
            keywords: keywords
        });
        
        allData.push(newData);
        userDataIds[msg.sender].push(newId);
        typeToDataIds[dataType].push(newId);
        
        emit DataUploaded(newId, msg.sender, dataType, block.timestamp);
        
        return newId;
    }
    
    /**
     * @dev 获取特定ID的医疗数据
     * @param id 数据ID
     * @return 医疗数据的完整信息
     */
    function getData(uint256 id) public view returns (
        uint256,
        address,
        string memory,
        string memory,
        string memory,
        uint256,
        string memory
    ) {
        require(id < allData.length, "Data does not exist");
        Data memory data = allData[id];
        
        return (
            data.id,
            data.owner,
            data.dataHash,
            data.dataType,
            data.metadata,
            data.timestamp,
            data.keywords
        );
    }
    
    /**
     * @dev 获取用户拥有的所有数据ID
     * @param user 用户地址
     * @return 数据ID数组
     */
    function getUserDataIds(address user) public view returns (uint256[] memory) {
        return userDataIds[user];
    }
    
    /**
     * @dev 获取特定类型的所有数据ID
     * @param dataType 数据类型
     * @return 数据ID数组
     */
    function getDataIdsByType(string memory dataType) public view returns (uint256[] memory) {
        return typeToDataIds[dataType];
    }
    
    /**
     * @dev 获取数据总数
     * @return 数据总数
     */
    function getDataCount() public view returns (uint256) {
        return allData.length;
    }
}