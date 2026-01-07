// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MyNFT
 * @dev ERC721 NFT 合约，支持铸造和 URI 设置
 */
contract MyNFT is ERC721URIStorage, Ownable {
    uint256 private _tokenIds;
    
    constructor(string memory name, string memory symbol) ERC721(name, symbol) Ownable(msg.sender) {}
    
    /**
     * @dev 铸造新的 NFT
     * @param to NFT 接收者地址
     * @param tokenURI NFT 元数据 URI
     * @return tokenId 新铸造的 NFT tokenId
     */
    function mint(address to, string memory tokenURI) public onlyOwner returns (uint256) {
        unchecked {
            _tokenIds++;
        }
        uint256 newTokenId = _tokenIds;
        
        _safeMint(to, newTokenId);
        _setTokenURI(newTokenId, tokenURI);
        
        return newTokenId;
    }
    
    /**
     * @dev 设置 NFT 的 URI
     * @param tokenId NFT tokenId
     * @param tokenURI NFT 元数据 URI
     */
    function setTokenURI(uint256 tokenId, string memory tokenURI) public onlyOwner {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        _setTokenURI(tokenId, tokenURI);
    }
    
    /**
     * @dev 获取当前总供应量
     * @return 总供应量
     */
    function totalSupply() public view returns (uint256) {
        return _tokenIds;
    }
}

