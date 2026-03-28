#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
DICOM Doctor v2.7.0 - 宿主 AI 自动分批处理模式

当没有 OpenAI API Key 时，自动切换到宿主 AI 模式：
1. 生成标准化的分批阅片请求
2. 提示宿主 AI 逐批读取图片并分析
3. 自动回填结果到 batch_XXX.filled.json
4. 每批完成后自动合并到总表 review_results.json
5. 支持断点续跑

作者: AI Assistant
版本: 2.7.0
日期: 2026-03-28
"""

import argparse
import json
import logging
import os
import sys
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Optional, Any

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger("dicom-doctor.host-ai-review")


class HostAIReviewer:
    """
    宿主 AI 分批阅片协调器
    
    此类负责：
    1. 将阅片任务分批输出给宿主 AI
    2. 接收宿主 AI 的分析结果
    3. 自动保存进度和合并结果
    """
    
    VERSION = "2.7.0"
    
    def __init__(self, manifest_path: str, output_dir: str, batch_size: int = 15):
        self.manifest_path = Path(manifest_path).resolve()
        self.output_dir = Path(output_dir).resolve()
        self.batch_size = batch_size
        
        # 加载 manifest
        with open(self.manifest_path, "r", encoding="utf-8") as f:
            self.manifest = json.load(f)
        
        self.total_images = self.manifest.get("total_images", 0)
        self.total_batches = self.manifest.get("total_batches", 0)
        self.requests = self.manifest.get("requests", [])
        
        # 设置目录
        self.batch_template_dir = self.output_dir / "review_batch_templates"
        self.batch_filled_dir = self.output_dir / "review_batch_filled"
        self.batch_filled_dir.mkdir(parents=True, exist_ok=True)
        
        # 检查已完成的批次
        self.completed_batches = self._get_completed_batches()
        
        logger.info(f"宿主 AI 阅片协调器 v{self.VERSION} 初始化完成")
        logger.info(f"总影像数: {self.total_images}, 总批次数: {self.total_batches}")
        logger.info(f"已完成批次: {len(self.completed_batches)}/{self.total_batches}")
    
    def _get_completed_batches(self) -> set:
        """获取已完成的批次索引"""
        completed = set()
        if self.batch_filled_dir.exists():
            for f in self.batch_filled_dir.glob("batch_*.filled.json"):
                try:
                    batch_num = int(f.stem.split("_")[1])
                    completed.add(batch_num)
                except (ValueError, IndexError):
                    pass
        return completed
    
    def get_next_batch(self) -> Optional[Dict]:
        """获取下一个待处理的批次"""
        for batch_index in range(1, self.total_batches + 1):
            if batch_index not in self.completed_batches:
                return self._load_batch(batch_index)
        return None
    
    def _load_batch(self, batch_index: int) -> Dict:
        """加载指定批次的模板"""
        batch_file = self.batch_template_dir / f"batch_{batch_index:03d}.json"
        if not batch_file.exists():
            raise FileNotFoundError(f"批次文件不存在: {batch_file}")
        
        with open(batch_file, "r", encoding="utf-8") as f:
            batch_data = json.load(f)
        
        return {
            "batch_index": batch_index,
            "data": batch_data,
            "items": batch_data.get("items", [])
        }
    
    def generate_host_ai_prompt(self, batch: Dict) -> str:
        """生成给宿主 AI 的阅片提示"""
        batch_index = batch["batch_index"]
        items = batch["items"]
        
        prompt_lines = [
            f"=" * 60,
            f"【DICOM Doctor 宿主 AI 阅片请求 - 批次 {batch_index}/{self.total_batches}】",
            f"=" * 60,
            f"",
            f"请逐张检视以下 {len(items)} 张胸部 CT 切片，并返回 JSON 格式的分析结果。",
            f"",
            f"【阅片规范】",
            f"1. 检视顺序：GGO 窗 → 肺窗 → 纵隔窗",
            f"2. ≥2mm 圆形/类圆形高密度影 → 疑似结节（禁止归为血管断面）",
            f"3. 扫及区域：甲状腺/肝脏上段/肾上腺（层面相关时检查）",
            f"",
            f"【回复 JSON 格式】",
            f"```json",
            f"{{",
            f'    "conclusion": "正常/异常/无法识别",',
            f'    "abnormality_desc": "异常描述（如有）",',
            f'    "confidence": "高/中/低",',
            f'    "details": "简要所见",',
            f'    "location": "精确到亚段+层面，如右肺中叶内段S5 (第50层)",',
            f'    "size_mm": "如3x2",',
            f'    "lung_rads": "如2类",',
            f'    "recommendation": "随访建议",',
            f'    "bounding_boxes": []',
            f"}}",
            f"```",
            f"",
            f"【待阅片列表】",
        ]
        
        for item in items:
            global_index = item.get("global_index", "?")
            slice_index = item.get("slice_index", "?")
            png_path = item.get("png_path", "")
            ggo_path = item.get("ggo_path", "")
            
            prompt_lines.append(f"")
            prompt_lines.append(f"--- 第 {global_index} 张 ---")
            prompt_lines.append(f"层面: {slice_index}")
            prompt_lines.append(f"肺窗: {png_path}")
            if ggo_path:
                prompt_lines.append(f"GGO窗: {ggo_path}")
            prompt_lines.append(f"Prompt: {item.get('prompt', '')[:100]}...")
        
        prompt_lines.append(f"")
        prompt_lines.append(f"=" * 60)
        prompt_lines.append(f"请读取以上图片，逐张分析后返回结果。")
        prompt_lines.append(f"=" * 60)
        
        return "\n".join(prompt_lines)
    
    def save_batch_results(self, batch_index: int, results: List[Dict]) -> str:
        """保存批次结果"""
        # 加载原始批次模板
        batch_file = self.batch_template_dir / f"batch_{batch_index:03d}.json"
        with open(batch_file, "r", encoding="utf-8") as f:
            batch_data = json.load(f)
        
        # 回填结果
        items = batch_data.get("items", [])
        for i, item in enumerate(items):
            if i < len(results):
                item["result"] = results[i]
        
        # 添加元数据
        batch_data["host_ai_review"] = {
            "generated_at": datetime.now().isoformat(timespec="seconds"),
            "version": self.VERSION,
            "reviewed_count": len([r for r in results if r.get("conclusion") in ["正常", "异常"]]),
        }
        
        # 保存到 filled 目录
        filled_file = self.batch_filled_dir / f"batch_{batch_index:03d}.filled.json"
        with open(filled_file, "w", encoding="utf-8") as f:
            json.dump(batch_data, f, ensure_ascii=False, indent=2)
        
        logger.info(f"批次 {batch_index} 结果已保存: {filled_file}")
        return str(filled_file)
    
    def merge_to_results(self, batch_index: int) -> str:
        """合并批次结果到总表"""
        # 调用 apply_review_batch.py 进行合并
        import subprocess
        
        filled_file = self.batch_filled_dir / f"batch_{batch_index:03d}.filled.json"
        results_file = self.output_dir / "review_results.json"
        stub_file = self.output_dir / "review_results_stub.json"
        
        # 确定输入的总表文件
        input_results = str(results_file) if results_file.exists() else str(stub_file)
        
        cmd = [
            sys.executable,
            str(Path(__file__).parent / "apply_review_batch.py"),
            "--manifest", str(self.manifest_path),
            "--results", input_results,
            "--batch-json", str(filled_file),
            "--output", str(results_file),
        ]
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
            if result.returncode == 0:
                logger.info(f"批次 {batch_index} 已合并到总表: {results_file}")
                return str(results_file)
            else:
                logger.error(f"合并批次 {batch_index} 失败: {result.stderr}")
                raise RuntimeError(f"合并失败: {result.stderr}")
        except Exception as e:
            logger.error(f"合并批次 {batch_index} 时出错: {e}")
            raise
    
    def get_progress(self) -> Dict:
        """获取当前进度"""
        return {
            "total_batches": self.total_batches,
            "completed_batches": len(self.completed_batches),
            "remaining_batches": self.total_batches - len(self.completed_batches),
            "completed_indices": sorted(self.completed_batches),
            "next_batch": self._get_next_batch_index(),
        }
    
    def _get_next_batch_index(self) -> Optional[int]:
        """获取下一个待处理的批次索引"""
        for i in range(1, self.total_batches + 1):
            if i not in self.completed_batches:
                return i
        return None
    
    def run_interactive(self):
        """交互式运行 - 输出提示供宿主 AI 处理"""
        print("\n" + "=" * 60)
        print("【DICOM Doctor 宿主 AI 分批阅片模式】")
        print("=" * 60)
        print(f"\n总批次数: {self.total_batches}")
        print(f"已完成: {len(self.completed_batches)}")
        print(f"待处理: {self.total_batches - len(self.completed_batches)}")
        print("\n请按以下步骤操作:")
        print("1. 我会逐批输出阅片请求")
        print("2. 你读取对应的 PNG 图片并分析")
        print("3. 返回 JSON 格式的结果")
        print("4. 我自动保存并合并到总表")
        print("=" * 60 + "\n")
        
        while True:
            batch = self.get_next_batch()
            if batch is None:
                print("\n✅ 所有批次已完成!")
                break
            
            batch_index = batch["batch_index"]
            print(f"\n{'='*60}")
            print(f"【批次 {batch_index}/{self.total_batches}】")
            print(f"{'='*60}\n")
            
            # 输出阅片提示
            prompt = self.generate_host_ai_prompt(batch)
            print(prompt)
            
            # 等待宿主 AI 返回结果
            print(f"\n{'='*60}")
            print("请返回上述图片的阅片结果（JSON 数组格式）:")
            print(f"{'='*60}\n")
            
            # 在实际实现中，这里需要与宿主 AI 交互
            # 目前输出提示，由宿主 AI 手动处理
            print("(等待宿主 AI 返回结果...)")
            break  # 每次只处理一个批次，避免输出过长


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description="DICOM Doctor v2.7.0 - 宿主 AI 自动分批阅片模式"
    )
    parser.add_argument(
        "--manifest",
        required=True,
        help="review_manifest.json 路径",
    )
    parser.add_argument(
        "--output",
        required=True,
        help="输出目录（包含 review_batch_templates/）",
    )
    parser.add_argument(
        "--batch-size",
        type=int,
        default=15,
        help="每批处理的图片数量（默认15）",
    )
    parser.add_argument(
        "--status",
        action="store_true",
        help="仅显示当前进度状态",
    )
    return parser.parse_args()


def main():
    """命令行入口"""
    args = parse_args()
    
    # 初始化协调器
    reviewer = HostAIReviewer(
        manifest_path=args.manifest,
        output_dir=args.output,
        batch_size=args.batch_size,
    )
    
    # 显示状态或运行交互模式
    if args.status:
        progress = reviewer.get_progress()
        print("\n【当前进度】")
        print(f"  总批次: {progress['total_batches']}")
        print(f"  已完成: {progress['completed_batches']}")
        print(f"  待处理: {progress['remaining_batches']}")
        if progress['next_batch']:
            print(f"  下一批: 批次 {progress['next_batch']}")
        else:
            print(f"  状态: ✅ 全部完成")
        return 0
    
    # 运行交互模式
    reviewer.run_interactive()
    return 0


if __name__ == "__main__":
    sys.exit(main())
