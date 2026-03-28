#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DICOM Doctor v2.8.0 - 宿主 AI 能力检测模块

检测当前宿主 AI 是否支持：
1. 图片识别（多模态）
2. 长文本处理
3. JSON 格式输出

作者: AI Assistant
版本: 2.8.0
日期: 2026-03-28
"""

import os
import sys
from typing import Dict, Any


class ModelCapabilityDetector:
    """宿主 AI 能力检测器"""
    
    # 已知支持图片识别的模型列表（部分）
    VISION_CAPABLE_MODELS = [
        "gpt-4", "gpt-4o", "gpt-4-turbo", "gpt-4-vision",
        "claude-3", "claude-3-opus", "claude-3-sonnet", "claude-3-haiku",
        "gemini", "gemini-pro", "gemini-ultra", "gemini-1.5",
        "qwen-vl", "qwen-vl-max", "qwen-vl-plus", "qwen2-vl",
        "llava", "yi-vl", "internvl", "cogvlm",
    ]
    
    # 已知不支持图片识别的模型列表
    TEXT_ONLY_MODELS = [
        "glm-4", "glm-5", "glm-5.0", "chatglm",
        "llama-2", "llama-3", "llama3",
        "qwen", "qwen2", "qwen1.5", "qwen-turbo",  # 纯文本版本
        "baichuan", "chatglm3", "chatglm2",
    ]
    
    @classmethod
    def detect_from_env(cls) -> Dict[str, Any]:
        """
        从环境变量检测模型能力
        
        Returns:
            Dict 包含:
                - model_name: 检测到的模型名称
                - supports_vision: 是否支持图片识别
                - supports_long_context: 是否支持长上下文
                - confidence: 检测置信度 (high/medium/low)
                - note: 说明信息
        """
        # 尝试从各种环境变量获取模型名称
        model_name = None
        env_vars = [
            "MODEL_NAME",
            "MODEL_ID", 
            "AI_MODEL",
            "LLM_MODEL",
            "OPENAI_MODEL",
            "GLM_MODEL",
            "ANTHROPIC_MODEL",
        ]
        
        for var in env_vars:
            if os.environ.get(var):
                model_name = os.environ.get(var).lower()
                break
        
        # 如果没有找到，尝试从其他线索推断
        if not model_name:
            model_name = cls._infer_from_context()
        
        if not model_name:
            return {
                "model_name": "unknown",
                "supports_vision": False,
                "supports_long_context": False,
                "confidence": "low",
                "note": "无法检测到模型信息，默认按纯文本模型处理",
            }
        
        # 检测是否支持视觉
        supports_vision = cls._check_vision_support(model_name)
        supports_long_context = cls._check_long_context_support(model_name)
        
        return {
            "model_name": model_name,
            "supports_vision": supports_vision,
            "supports_long_context": supports_long_context,
            "confidence": "high" if supports_vision else "medium",
            "note": cls._generate_note(model_name, supports_vision),
        }
    
    @classmethod
    def _infer_from_context(cls) -> str:
        """从上下文推断模型类型"""
        # 检查是否有特定的视觉模型标志
        if os.environ.get("VISION_ENABLED") == "true":
            return "vision-capable-model"
        
        # 检查是否有 GLM 相关环境变量
        if os.environ.get("GLM_API_KEY") or os.environ.get("ZHIPU_API_KEY"):
            return "glm-model"
        
        # 检查运行的进程或库
        try:
            import importlib.util
            # 检查是否导入了多模态相关的库
            vision_libs = ["pillow", "cv2", "opencv-python"]
            for lib in vision_libs:
                if importlib.util.find_spec(lib):
                    return "potentially-vision-capable"
        except:
            pass
        
        return None
    
    @classmethod
    def _check_vision_support(cls, model_name: str) -> bool:
        """检查模型是否支持视觉"""
        model_lower = model_name.lower()
        
        # 检查是否在已知支持列表中
        for vision_model in cls.VISION_CAPABLE_MODELS:
            if vision_model in model_lower:
                return True
        
        # 检查是否在已知不支持列表中
        for text_model in cls.TEXT_ONLY_MODELS:
            if text_model in model_lower:
                return False
        
        # 默认保守处理：假设不支持
        return False
    
    @classmethod
    def _check_long_context_support(cls, model_name: str) -> bool:
        """检查模型是否支持长上下文"""
        # 简化检测：假设现代模型都支持至少 4K 上下文
        long_context_models = [
            "gpt-4", "claude-3", "gemini", "glm-4", "glm-5",
            "qwen", "llama-3", "yi", "baichuan2",
        ]
        
        model_lower = model_name.lower()
        for model in long_context_models:
            if model in model_lower:
                return True
        
        return False
    
    @classmethod
    def _generate_note(cls, model_name: str, supports_vision: bool) -> str:
        """生成说明信息"""
        if supports_vision:
            return f"检测到模型 {model_name} 支持图片识别，可以自动完成阅片"
        else:
            return f"检测到模型 {model_name} 不支持图片识别，需要使用外部 API 或手动处理"


def detect_host_ai_capabilities() -> Dict[str, Any]:
    """
    便捷的检测函数
    
    Returns:
        模型能力信息字典
    """
    detector = ModelCapabilityDetector()
    return detector.detect_from_env()


def print_capability_report() -> Dict[str, Any]:
    """打印能力检测报告"""
    capabilities = detect_host_ai_capabilities()
    
    print("\n" + "=" * 60)
    print("【宿主 AI 能力检测报告】")
    print("=" * 60)
    print(f"模型名称: {capabilities['model_name']}")
    print(f"图片识别: {'✅ 支持' if capabilities['supports_vision'] else '❌ 不支持'}")
    print(f"长上下文: {'✅ 支持' if capabilities['supports_long_context'] else '❌ 不支持'}")
    print(f"检测置信度: {capabilities['confidence']}")
    print(f"\n说明: {capabilities['note']}")
    print("=" * 60 + "\n")
    
    return capabilities


def get_recommended_action(capabilities: Dict[str, Any]) -> str:
    """
    根据模型能力返回推荐的操作方案
    
    Returns:
        推荐的操作说明字符串
    """
    if capabilities['supports_vision']:
        return """
推荐操作方案：
✅ 当前模型支持图片识别，可以使用宿主 AI 自动阅片模式

运行命令：
  python scripts/host_ai_review.py \\
    --manifest <manifest路径> \\
    --output <输出目录> \\
    --auto-continue

系统会自动连续处理所有批次，无需人工干预。
"""
    else:
        return """
推荐操作方案（按优先级排序）：

1️⃣  使用 OpenAI API（推荐，最稳定）
   运行命令：
   python scripts/main.py --input <dicom文件> --output <输出目录> \\
     --auto-review-model gpt-4o \\
     --auto-review-api-key <你的API Key>

2️⃣  使用 Claude API
   运行命令：
   python scripts/main.py --input <dicom文件> --output <输出目录> \\
     --auto-review-model claude-3-opus-20240229 \\
     --auto-review-api-key <你的API Key>

3️⃣  手动分批处理
   如果无法使用外部 API，可以手动将图片发送给支持视觉的 AI 进行分析

💡 获取 API Key：
   - OpenAI: https://platform.openai.com/api-keys
   - Anthropic: https://console.anthropic.com/
"""


if __name__ == "__main__":
    caps = print_capability_report()
    print(get_recommended_action(caps))
