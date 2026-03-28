#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DICOM Doctor v2.9.0 - 宿主 AI 能力检测模块

检测当前宿主 AI 是否支持：
1. 图片识别（多模态）
2. 长文本处理
3. JSON 格式输出

设计原则（v2.9.0 重要变更）：
  - SKILL.md 已明确要求使用多模态视觉模型，因此当无法确定模型能力时，
    **默认假定支持视觉**（乐观策略），而非保守地假定不支持。
  - 保守策略会导致宿主 AI 被错误引导到"纯文本模式"，进而放弃阅片，
    这比乐观策略的风险（尝试阅片后发现不支持再降级）更严重。
  - 新增 detect_from_model_name() 方法，支持从 --model-name 参数直接检测。

作者: AI Assistant
版本: 2.9.0
日期: 2026-03-28
"""

import os
import sys
from typing import Dict, Any


class ModelCapabilityDetector:
    """宿主 AI 能力检测器"""
    
    # 已知支持图片识别的模型列表（持续更新）
    # 匹配规则：model_name.lower() 中包含列表中任一字符串即视为匹配
    VISION_CAPABLE_MODELS = [
        # OpenAI 系列
        "gpt-4", "gpt-4o", "gpt-4-turbo", "gpt-4-vision",
        "gpt-5", "gpt-5.4", "gpt-5o",
        "o1", "o3", "o4-mini",
        # Anthropic Claude 系列
        "claude-3", "claude-3-opus", "claude-3-sonnet", "claude-3-haiku",
        "claude-3.5", "claude-3.7", "claude-4",
        "claude-4.6-opus", "claude-4-opus", "claude-4-sonnet",
        # Google Gemini 系列
        "gemini", "gemini-pro", "gemini-ultra", "gemini-1.5",
        "gemini-2", "gemini-2.5",
        # Kimi / Moonshot 系列
        "kimi", "kimi-k2", "kimi-k2.5", "moonshot-v",
        # 阿里通义千问 VL 系列
        "qwen-vl", "qwen-vl-max", "qwen-vl-plus", "qwen2-vl",
        "qwen2.5-vl", "qwen-omni",
        # 其他多模态模型
        "llava", "yi-vl", "yi-vision", "internvl", "cogvlm",
        "glm-4v", "glm-5v",  # 智谱视觉版本
        "step-1v", "step-2v",  # 阶跃星辰
        "doubao-vision", "doubao-1.5-vision",  # 豆包视觉版
        "hunyuan-vision",  # 腾讯混元视觉版
        "ernie-4-turbo",  # 百度文心视觉版
        "deepseek-vl",  # DeepSeek 视觉版
    ]
    
    # 已知不支持图片识别的纯文本模型列表
    # 注意：只有明确确认为纯文本的模型才放在这里
    TEXT_ONLY_MODELS = [
        "chatglm", "chatglm2", "chatglm3",
        "llama-2", "llama-3", "llama3",
        "qwen-turbo", "qwen-plus", "qwen-max",  # 纯文本版本（不含 -vl 后缀）
        "baichuan", "baichuan2",
        "deepseek-chat", "deepseek-coder",  # DeepSeek 纯文本版
    ]
    
    @classmethod
    def detect_from_model_name(cls, model_name: str) -> Dict[str, Any]:
        """
        从显式传入的模型名称检测能力（推荐方式）。
        
        当宿主 AI 通过 --model-name 参数传入模型名称时使用此方法，
        比从环境变量推断更可靠。
        
        Args:
            model_name: 模型名称（如 "claude-4.6-opus"、"GPT-5.4"）
        
        Returns:
            模型能力信息字典
        """
        if not model_name:
            return cls.detect_from_env()
        
        model_lower = model_name.lower().strip()
        supports_vision = cls._check_vision_support(model_lower)
        supports_long_context = cls._check_long_context_support(model_lower)
        
        return {
            "model_name": model_name,
            "supports_vision": supports_vision,
            "supports_long_context": supports_long_context,
            "confidence": "high",
            "note": cls._generate_note(model_name, supports_vision),
            "detection_method": "model_name_param",
        }

    @classmethod
    def detect_from_env(cls) -> Dict[str, Any]:
        """
        从环境变量检测模型能力（回退方式）。
        
        v2.9.0 重要变更：当无法确定模型时，默认假定支持视觉。
        理由：SKILL.md 已明确要求使用多模态视觉模型，保守策略会导致
        宿主 AI 被错误引导到纯文本模式而放弃阅片，风险远大于乐观策略。
        
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
            # v2.9.0 变更：无法检测时默认假定支持视觉（乐观策略）
            return {
                "model_name": "unknown",
                "supports_vision": True,
                "supports_long_context": True,
                "confidence": "low",
                "note": (
                    "无法检测到具体模型信息。根据 SKILL.md 要求（必须使用多模态视觉模型），"
                    "默认假定当前模型支持视觉能力。如果实际不支持，阅片时会自动降级处理。"
                ),
                "detection_method": "default_optimistic",
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
            "detection_method": "env_var",
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
        """
        检查模型是否支持视觉。
        
        v2.9.0 变更：对于未知模型，默认假定支持视觉（乐观策略）。
        只有明确在 TEXT_ONLY_MODELS 列表中的模型才返回 False。
        """
        model_lower = model_name.lower()
        
        # 先检查是否在已知不支持列表中（精确排除）
        for text_model in cls.TEXT_ONLY_MODELS:
            if text_model in model_lower:
                # 但要排除误匹配：如 "qwen-vl" 包含 "qwen" 但实际支持视觉
                for vision_model in cls.VISION_CAPABLE_MODELS:
                    if vision_model in model_lower:
                        return True
                return False
        
        # 检查是否在已知支持列表中
        for vision_model in cls.VISION_CAPABLE_MODELS:
            if vision_model in model_lower:
                return True
        
        # v2.9.0 变更：未知模型默认假定支持视觉
        # 理由：SKILL.md 要求使用多模态模型，乐观策略比保守策略风险更低
        return True
    
    @classmethod
    def _check_long_context_support(cls, model_name: str) -> bool:
        """
        检查模型是否支持长上下文。
        
        v2.9.0: 2024-2026 年的主流模型基本都支持 128K+ 上下文，
        默认返回 True。
        """
        # 现代模型（2024+）基本都支持长上下文，默认 True
        return True
    
    @classmethod
    def _generate_note(cls, model_name: str, supports_vision: bool) -> str:
        """生成说明信息"""
        if supports_vision:
            return f"模型 {model_name} 支持图片识别，将使用宿主 AI 直接阅片"
        else:
            return (
                f"模型 {model_name} 被识别为纯文本模型，不支持图片识别。"
                f"建议使用外部视觉 API（--auto-review-model）或切换到多模态模型。"
                f"如果判断有误，宿主 AI 可忽略此检测结果，直接按流程阅片。"
            )


def detect_host_ai_capabilities(model_name: str = None) -> Dict[str, Any]:
    """
    便捷的检测函数。
    
    Args:
        model_name: 可选的模型名称（如通过 --model-name 传入）。
                    提供时优先使用，比环境变量检测更可靠。
    
    Returns:
        模型能力信息字典
    """
    detector = ModelCapabilityDetector()
    if model_name:
        return detector.detect_from_model_name(model_name)
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
