---
name: dicom-doctor
description: "Use this skill whenever the user mentions DICOM files, medical imaging, CT scans, MRI scans, or anything related to medical image conversion or chest CT analysis. This includes: converting any DICOM files (CT, MRI, PET, etc.) to PNG images, analyzing chest CT for lung nodules/masses/calcifications, generating PDF radiology review reports for chest CT, or any mention of 医学影像, 阅片, DICOM, 胸部CT, CT扫描, MRI, 核磁, 核磁共振, 肺结节, 影像报告, dcm文件, DICOM转PNG, 影像转换. Even if the user just says 'convert this DICOM to PNG' or 'look at this chest CT' or mentions a .dcm file path, activate this skill. Also trigger for requests about Lung-RADS classification, MIP reconstruction, super-resolution enhancement, or window settings (lung/mediastinum/GGO). Note: AI-assisted review (阅片) currently only supports chest CT; DICOM-to-PNG conversion supports all modalities including MRI, abdominal CT, brain MRI, etc."
---

# DICOM Doctor — AI 辅助医学影像阅片 Skill

**当前版本：v2.6.0**

## 🚨 铁律（违反任何一条即为严重失败）

> 以下规则对阅片质量和患者安全至关重要。宿主 AI 必须逐条遵守，不得以任何理由违反。

1. **逐层确认式全量阅片，零遗漏**：必须逐张检视**全部** PNG 切片，**每一层都必须输出明确结论**（正常/异常/无法识别），不得以任何方式跳过或简化。3mm 结节仅出现在约 3 层中，跳过即漏诊。**严格禁止**：抽样、跳过、拼图（Collage）、以任何理由降级为部分阅片。
2. **PNG 必须保存到 png/ 子目录**：运行 `main.py` 后，PNG 文件输出在 `<output_dir>/<时间戳>/png/` 下。必须验证该目录存在且包含文件。
3. **必须生成 PDF 报告（对标模板，不得省略）**：最终**必须**在输出目录中生成 `.pdf` 格式的医院风格影像检查报告。报告格式**必须严格对标** `references/AI_chest_CT_report_template.pdf` 模板，包含：检查信息表、AI 检视统计、检查所见（逐条异常发现）、异常影像详情（嵌入图片+描述+分级）、诊断意见、随访建议、分级参考表、免责声明。**如果最终输出目录中不存在 .pdf 文件，则本次任务判定为失败。**
4. **立即执行，不得确认或中途停止**：收到用户指令后，**直接开始执行第 1 步（运行 main.py）**，禁止在开始前插入任何形式的确认环节。以下行为**全部禁止**：
   - "请确认以下信息""我先确认一下""没问题的话我就开始"
   - **"先补个最小确认""快速确认一下"**——这也是确认，同样禁止
   - 列出用户信息/环境信息后问"这些信息都对吗？"
   - 复述用户的输入路径、姓名、系统环境后请求用户确认
   - 中途暂停询问"是否继续""需要我继续吗"
   - **提供"方案A/方案B"让用户选择**——这本质上是变相暂停，同样禁止
   - **以"切片数量太多""窗位质量差""narrow_ggo 过曝"等为由暂停或降级**——不管数据量有多大、某个窗位看起来多差，流程都必须跑完，问题记录到报告中即可
   
   **唯一例外**：输入文件路径确实不存在（文件系统报错），才可以询问用户。其他一切情况，直接开始干活。
5. **GGO 窗优先**：胸部 CT 阅片时，每批切片必须先独立检视**高灵敏度 GGO 窗**（`narrow_ggo/` 子目录），再检视 GGO 窗（`ggo/` 子目录），然后检视肺窗和纵隔窗。极淡的纯磨玻璃结节可能**只有在高灵敏度 GGO 窗**下才能看到。
   
   **narrow_ggo 过曝/质量差时的自动降级策略**：如果第1批 narrow_ggo 图片出现过曝（几乎全白/全黑、对比度极低、无法辨别肺实质），则**自动将 GGO 主力检测窗降级为 `ggo/`**，后续批次跳过 narrow_ggo、以 ggo 窗为 GGO 检测首选。此降级**由 AI 自行判断执行，不得暂停询问用户**。降级事实须在最终报告"检查所见"中注明（如"因 narrow_ggo 窗过曝，本次 GGO 检测以 ggo 窗为主力"）。
6. **双侧扫查**：每张图片都要同时检查左肺和右肺。只在单侧发现异常时，必须回头复查对侧。
7. **CAD 预检候选必须验证，验证结论必须三选一**：如果 `main.py` 输出了 CAD 自动预检结果（`cad_candidates.json` 和 `cad_annotations/` 标注图），阅片时**必须逐一验证这些候选区域**。v2.2.0 CAD 使用基于5个真实结节(2例胸部CT)校准的7维评分排序（⭐>0.9 高度可疑，0.8-0.9 中度可疑），评分维度含球形度、elongation、大小、HU、z层数、密度一致性、血管排除。厚层CT(≥1mm)额外启用 binary_closing + 膨胀保护 + 密度峰值提取 + 大连通域2D子候选提取 + GGO逐层2D去噪，提升微小结节检出率。高分候选必须优先验证。CAD 同时输出全切面图和四窗位合成图，方便对照判断。
   
   **每个 CAD 候选验证后，结论必须三选一（禁止"看了但不报"）：**
   - ✅ **确认结节**→ conclusion="异常"，正常填写所有异常字段
   - ⚠️ **可疑/不确定**（如"形态更符合血管但不能排除结节"、"需随访"）→ **conclusion 必须填"异常"**，abnormality_desc 中如实描述可疑特征及不确定原因，confidence 填"低"或"中"
   - ❌ **明确排除**（如确认为正常血管分支、伪影等）→ conclusion="正常"，在 details 中注明排除原因
   
   > **核心原则：宁可多报一个可疑发现让医生复核，也不能吞掉一个可能的结节。CAD score≥0.8 的候选，除非有充分证据排除，否则必须标记为"异常"。**
   
8. **报告必须包含两个维度的完整异常数据（缺一不可）**：
   - **结节维度（聚合汇总）**：同一结节跨层面的异常记录必须合并为独立结节，输出结节汇总表（含位置、类型、大小、出现层面范围、分级、置信度）。这是临床医生快速评估的核心信息。
   - **层面维度（逐层明细）**：每个异常层面的详细数据必须完整保留（含层面序号、位置、异常描述、大小、分级）。这是精确定位复核的关键依据。
   - **❌ 严禁**只输出结节汇总而省略逐层明细。**❌ 严禁**只输出逐层明细而省略结节汇总。两者都必须完整输出到报告中。
9. **可疑/疑似/不确定的发现必须标记为"异常"并进入报告**：
   - 阅片时如果描述中出现了"可疑"、"疑似"、"不确定"、"不能排除"、"需随访"、"需谨慎对待"等措辞，该层面的 conclusion **必须填"异常"**，不得填"正常"
   - **❌ 严禁"描述了但不报"**：AI 在阅片描述中提到了异常特征（如"小圆点状高密度影"、"可疑微小结节"、"需随访"），但 conclusion 填了"正常"——这是最危险的漏诊模式，因为 AI 明明看到了却没让它进入报告
   - **❌ 严禁"自行排除"高分 CAD 候选**：对 CAD score≥0.8 的候选，AI 不得仅凭"更像血管"就标正常。除非有多维度明确证据（形态、密度、连续性全部指向血管），否则必须标异常
   - 报告中如实标注不确定程度即可（confidence 填"低"或"中"），**让最终判断权留给医生**
10. **反保守偏见（铁律级）**：≥2mm 的圆形/类圆形高密度影**必须优先报告为疑似结节**，禁止轻易归为"血管断面"。具体规则：
    - HU 值在实性结节范围（≥20 HU）且直径≥2mm → 默认报告为可疑结节，除非有充分血管证据
    - 连续2层以上出现的局灶性密度增高 → 必须报告
    - 孤立圆形高密度影 → 即使只在1层出现也必须报告（可标注"仅单层可见，需随访确认"）
    - **"更像血管"不等于"是血管"**——不确定就报异常，让医生判断

## 概述

DICOM Doctor 是一个 AI 辅助医学影像处理 skill，提供两项核心能力：

- **DICOM 转 PNG**：支持所有影像类型（胸部CT、腹部CT、头颅MRI、腹部MRI、核磁等），自动选择合适的窗位策略
- **AI 辅助阅片**：**当前版本仅支持胸部CT**，逐张检视全部切片并生成 PDF 检查报告

工作流程：接收 DICOM 文件或 ZIP 压缩包 → 自动识别影像类型 → 转换为 PNG → （胸部CT）AI 逐张检视全部影像 → 生成 PDF 检查报告

### DICOM 转 PNG 支持的影像类型

| 影像类型 | 自动识别 | 窗位策略 | MIP | GGO窗 |
|----------|---------|---------|-----|-------|
| **胸部CT** | ✅ Modality=CT + BodyPart=CHEST | 肺窗/纵隔窗/GGO窗/骨窗 | ✅ | ✅ |
| **腹部CT** | ✅ Modality=CT + BodyPart=ABDOMEN | 软组织窗/肝脏窗/骨窗 | ❌ | ❌ |
| **头颅MRI** | ✅ Modality=MR + BodyPart=HEAD/BRAIN | T1/T2/FLAIR/DWI 自适应 | ❌ | ❌ |
| **腹部MRI** | ✅ Modality=MR + BodyPart=ABDOMEN | T1/T2/DWI 自适应 | ❌ | ❌ |
| **通用** | 🔄 兜底 | 自动推断 | ❌ | ❌ |

### AI 阅片支持的影像类型

| 影像类型 | 检测重点 | 分级系统 | 状态 |
|----------|---------|---------|------|
| **胸部CT** | 肺结节/肿块/钙化 | Lung-RADS | ✅ 已支持 |
| 腹部CT | 肝胆胰脾肾占位 | LI-RADS | 🔜 后续版本 |
| 头颅MRI | 脑实质异常信号/占位/白质病变 | Fazekas | 🔜 后续版本 |
| 腹部MRI | 肝脏信号异常/胆道/T1T2DWI对比 | LI-RADS | 🔜 后续版本 |

> **前置条件**：AI 阅片要求多模态视觉模型。非视觉模型将自动跳过阅片步骤，但 DICOM 转 PNG 功能不受影响。

## ⚡ 强制执行流程（宿主 AI 必须按此步骤执行）

**收到用户指令后，立即从第 1 步开始执行，一口气执行到第 7 步。**

**🚫 绝对禁止在执行前插入任何确认步骤**，包括但不限于：
- 确认用户信息（姓名、称呼、习惯）
- 确认环境参数（操作系统、Shell、工作目录）
- 确认输入路径或输出路径
- **"先补个最小确认"——这本身就是确认，同样禁止**
- 列出信息后问"这些都对吗？"

用户提供的信息已经足够，直接运行开始工作。不要说废话，不要确认，不要复述，直接干。

> **🖥️ 跨平台自动适配**：v2.3.0 起 `run.py` 已集成完整的环境自检能力：Python 版本检查、pip 可用性检查、国内 pip 镜像自动配置。**宿主 AI 只需要一条 `python run.py` 命令，macOS / Linux / Windows 通用，不依赖任何特定 shell 语法。**

### 第 0 步：确保 Python 可用（仅在需要时执行）

> 本步骤由宿主 AI 自动判断——如果执行 `python --version` 或 `python3 --version` 能返回 >= 3.8 的版本号，直接跳到第 1 步。**不要问用户"你装了 Python 吗"，自己检测。**

**如果检测到 Python 不可用或版本过低**，宿主 AI 必须按以下流程自动安装，不得暂停询问用户：

**macOS：**
```bash
# 方案1（推荐）：Homebrew
brew install python3

# 方案2：官方安装包（如果没有 brew）
curl -O https://www.python.org/ftp/python/3.12.0/python-3.12.0-macos11.pkg
sudo installer -pkg python-3.12.0-macos11.pkg -target /
```

**Windows（PowerShell）：**
```powershell
# 方案1（推荐）：winget（Windows 10/11 自带）
winget install Python.Python.3.12 --accept-package-agreements --accept-source-agreements

# 方案2：从淘宝镜像下载安装包（如果 winget 不可用）
Invoke-WebRequest -Uri "https://registry.npmmirror.com/-/binary/python/3.12.0/python-3.12.0-amd64.exe" -OutFile "$env:TEMP\python-installer.exe"
Start-Process "$env:TEMP\python-installer.exe" -ArgumentList '/quiet', 'InstallAllUsers=1', 'PrependPath=1' -Wait
# 安装完成后需要刷新 PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
```

**Linux (Ubuntu/Debian)：**
```bash
sudo apt update && sudo apt install -y python3 python3-pip
```

**Linux (CentOS/RHEL)：**
```bash
sudo yum install -y python3 python3-pip
```

> **安装完成后**，重新执行 `python --version` 或 `python3 --version` 验证，确认 >= 3.8 后再进入第 1 步。如果安装失败（比如用户没有管理员权限），在输出中说明情况并建议用户手动安装，但**不要反复询问**。

### 第 1 步：运行 run.py（自动创建输出目录并启动 main.py）

```
python <skill_path>/scripts/run.py --input <用户提供的DICOM路径> --workspace <用户工作区> --model-name <当前模型名称> --strict-review
```

> **说明**：
> - `run.py` 会自动在 `<workspace>` 下创建 `dicom_output_<时间戳>/` 输出目录，并将 `--input`、`--output` 及所有额外参数透传给 `main.py`
> - **macOS/Linux** 上 `python` 和 `python3` 均可；**Windows** 上通常用 `python`。`run.py` 内部会自动选择正确的解释器
> - 如需直接指定输出目录，可用 `--output <路径>` 替代 `--workspace`
>
> **交付/正式阅片场景必须启用 `--strict-review`**：这样当逐张阅片结果仍有"待检视/无法识别"时，流程会直接拦截，不会误生成正式报告。
>
> **v1.3.0 起可选直连外部视觉模型自动回填**：若你手头有 OpenAI 兼容多模态接口，可追加 `--auto-review-model --auto-review-api-base --auto-review-api-key(或环境变量)`，让批次模板自动逐批回填并直接合并出 `review_results.json`。

### 第 2 步：验证 PNG 输出

`run.py` 执行完成后会打印输出目录路径。**必须检查** PNG 文件已正确输出——列出 `<输出目录>/<时间戳>/png/` 下的子目录和文件数量。胸部CT应有 `lung/`、`mediastinum/`、`ggo/`、`narrow_ggo/` 四个子目录。

如果 PNG 文件数为 0，说明转换失败，需要排查日志。**不得跳过此验证直接进入阅片。**

### 第 3 步：全量阅片（核心步骤，耗时最长）

**这是最关键的步骤。必须逐层确认式全量阅片——逐张检视全部 PNG 图片，一张不落，每张都必须输出明确结论。**

> **v2.4.7 上下文优化**：`main.py` 的 `review_manifest.json` 中，每层的 prompt 和 CAD hint 已自动优化：
> - **📌重点层**（CAD 候选 ±5 层范围）：使用完整版 prompt + 仅本层附近的 CAD 候选提示
> - **⚡快扫层**（远离 CAD 候选）：使用精简版 prompt（约 600 字 vs 完整版 5000 字）
> - 这使得 prompt 总 token 消耗降低约 50-70%，大幅缓解上下文窗口耗尽问题
> - **宿主 AI 仍需逐张检视全部图片**，快扫层只是 prompt 更简洁，检视精度要求不变

阅片流程：
1. 列出所有需要检视的 PNG 文件（按 lung/ 子目录中的文件列表）
2. 计算总数 N，按每批 15-20 张分批。**必须在开始时公布总数 N 和预计批次数**
3. 对每个批次：
   a. **首先**独立检视该批次的高灵敏度 GGO 窗图片（`narrow_ggo/` 子目录中同名文件）——检测极淡的纯磨玻璃结节
   b. **然后**独立检视该批次的 GGO 窗图片（`ggo/` 子目录中同名文件）——寻找磨玻璃结节
   c. **再**检视肺窗图片（`lung/` 子目录）——寻找实性结节、肿块
   d. **同时**参考纵隔窗图片（`mediastinum/` 子目录）——验证结节密度、检查淋巴结
   e. **为每张图片输出 JSON 格式的分析结论**（结论字段不能为空）
   f. 批次完成后**报告进度**："已完成第 M/N 批，累计检视 X/Y 张"
4. **必须循环直到最后一张图片都检视完毕**，然后进入汇总

**阅片红线（违反即严重医疗失误，任务判定失败）：**
- ❌ **严禁抽样**：禁止只看几张代表性切片
- ❌ **严禁拼图**：禁止将多张切片缩小拼接成一张（Collage）
- ❌ **严禁跳层**：禁止以"层面正常""已到达腹部区域"等理由跳过剩余层面
- ❌ 禁止中途停下来问"需要我继续检视剩余图片吗？"
- ❌ **严禁方案选择式暂停**：禁止在阅片过程中向用户提供"方案A/方案B"或"两个选项"让用户决策。遇到任何技术问题（窗位质量差、数据量大、图片过曝等），AI 必须自行按降级策略处理并继续，不得暂停
- ❌ **严禁以数据量为由降级**：无论切片数量是 100 还是 1000+、窗位组合产生多少张图片，都必须按流程全量跑完。担心"会话太长""token 不够"不是停下来的理由
- ✅ 必须每批 15-20 张，循环直到最后一张
- ✅ 每张图片都要输出明确结论（正常/异常/无法识别）
- ✅ 遇到窗位质量问题时，自行降级处理并在报告中注明

> 📖 详细的三阶段阅片策略和禁止行为清单，参见 `references/review_strategy.md`

### 第 4 步：汇总全部阅片发现

全部切片检视完毕后，汇总所有发现：
- 统计正常/异常/无法识别的层面数
- 对重复出现的结节进行**跨层面去重合并**
- 为每个确认的结节/病灶给出分级：Lung-RADS（胸部CT）/ LI-RADS（腹部CT/MRI）/ Fazekas（头颅MRI白质病变）
- 检查扫及区域（甲状腺、肝脏上段、肾上腺）是否已覆盖

**🚨 报告必须同时包含两个维度的完整数据（缺一不可，违反即为报告不合格）：**

1. **结节维度（聚合汇总）**：同一结节跨层面的异常记录必须合并为独立结节，每个结节必须包含：
   - 位置（精确到亚段）
   - 类型（GGO/实性/部分实性）
   - 大小（长径×短径 mm）
   - 出现层面范围（第 NNN-NNN 层）
   - 出现层面数
   - 分级（Lung-RADS/LI-RADS/Fazekas）
   - 置信度
   
2. **层面维度（逐层明细）**：每个异常层面的详细数据必须完整保留，包含：
   - 层面序号（如第 218/292 层）
   - 位置描述
   - 异常描述
   - 大小
   - 分级

> **❌ 严禁**只输出结节汇总而省略逐层明细——逐层明细是定位复核的关键依据。
> **❌ 严禁**只输出逐层明细而省略结节汇总——结节汇总是临床医生快速评估的核心信息。
> **✅ 两者都必须完整输出**，先结节汇总表，再逐层异常明细列表。

### 第 5 步：保存阅片结果到 JSON 文件

全量阅片完成后，将所有阅片结果保存为 JSON 文件，格式如下：

```bash
# 在输出目录的时间戳子目录下保存
REVIEW_RESULTS_JSON="$OUTPUT_DIR/<时间戳>/review_results.json"
```

> **v1.3.0 起闭环继续升级：** `main.py` 在未提供正式阅片结果时，会自动导出：
> - `review_requests.md`：逐张阅片请求与 prompt 汇总
> - `review_manifest.json`：结构化请求清单（后续会用于结果校验）
> - `review_results_stub.json`：待回填的占位结果 JSON
> - `review_batch_templates/batch_XXX.json`：按批拆好的回填模板
> - `review_batch_filled/`：使用外部视觉模型自动回填后生成的批次结果目录（按需生成）
>
> 推荐做法：
> 1. **自动模式**：直接运行 `main.py --auto-review-model ...`，或后续运行 `auto_review_batches.py`，让外部视觉模型逐批回填并自动合并出 `review_results.json`
> 2. **手工/半自动模式**：按批填写 `review_batch_templates/batch_XXX.json` 中每个 `item.result`
> 3. 每完成一批，运行 `apply_review_batch.py` 合并到总表 JSON
> 4. 全部批次完成后，再调用 `generate_report.py` 生成正式报告

JSON 格式（数组，每个元素对应一张切片的阅片结果）：
```json
[
  {
    "png_name": "IM-0001.png",
    "dicom_name": "IM-0001.dcm",
    "png_path": "/absolute/path/to/png/lung/IM-0001.png",
    "conclusion": "正常",
    "abnormality_desc": "",
    "confidence": "高",
    "details": "该层面显示双肺纹理清晰...",
    "location": "",
    "size_mm": "",
    "lung_rads": "",
    "recommendation": "",
    "slice_index": "1/100",
    "slice_location": "-120.5",
    "bounding_boxes": []
  },
  {
    "png_name": "IM-0050.png",
    "dicom_name": "IM-0050.dcm",
    "png_path": "/absolute/path/to/png/lung/IM-0050.png",
    "conclusion": "异常",
    "abnormality_desc": "右肺中叶内段(S5)可见约3mm×2mm实性结节",
    "confidence": "中",
    "details": "...",
    "location": "右肺中叶内段(S5) (第50层)",
    "size_mm": "3x2",
    "lung_rads": "2类",
    "recommendation": "建议12个月低剂量CT随访",
    "slice_index": "50/100",
    "slice_location": "-85.3",
    "bounding_boxes": [{"x": 0.35, "y": 0.42, "width": 0.05, "height": 0.04}]
  }
]
```

### 第 6 步：生成 PDF 报告

使用独立的报告生成脚本，从阅片结果 JSON 生成医院风格的 PDF 报告：

```bash
# 生成 PDF 和 Markdown 报告（默认正式模式，会校验 manifest 且拒绝待检视条目）
python3 <skill_path>/scripts/generate_report.py \
  --results "$REVIEW_RESULTS_JSON" \
  --manifest "$OUTPUT_DIR/<时间戳>/review_manifest.json" \
  --output "$OUTPUT_DIR/<时间戳>" \
  --input-path <原始DICOM路径> \
  --imaging-type chest_ct \
  --model-name <当前模型名称>
```

报告将参照 `references/AI_chest_CT_report_template.pdf` 的格式，包含：
- 检查信息表（检查类型、日期、影像数量、窗口类型等）
- AI 检视统计（总数、正常、异常、无法识别）
- **检查所见（🚨 必须同时包含两个维度，缺一不可）：**
  - **结节聚合汇总表**：将同一结节跨层面的异常记录自动合并为独立结节，展示结节总数、位置、类型、大小、出现层面范围和分级
  - **逐层异常明细**：每个异常层面的详细数据（位置、异常描述、大小、分级）——供精确定位复核
- 异常影像详情（嵌入照片 + 异常描述 + 分级分类）
- 诊断意见（按聚合结节/病灶输出，而非逐层重复） + 随访建议
- 分级参考表（Lung-RADS / LI-RADS / Fazekas，根据影像类型自动选择）
- 免责声明

验证报告已生成：
```bash
find "$OUTPUT_DIR" -name "*.pdf" -type f
find "$OUTPUT_DIR" -name "*.md" -type f
```

**⚠️ 如果 PDF 文件不存在或大小为 0，必须排查原因并重试，不得跳过。没有 PDF 报告 = 任务失败。**

> 💡 **备注**：`main.py` 运行时也会尝试生成 PDF 报告（阶段 4），但因为 AI 阅片结果需要宿主 AI 回填，
> 所以推荐使用 `generate_report.py` 在阅片完成后独立生成报告，确保报告内容完整。

### 第 7 步：向用户呈现结果

将以下信息返回给用户：
1. 输出目录路径
2. PNG 文件总数和目录结构
3. **全量阅片完成确认**：明确列出"共 N 张切片，已全部逐张检视完毕"，以及正常 X 张、异常 Y 张的统计
4. 每个异常结节的详细信息（位置、大小、分级）
5. **PDF 报告路径**（必须存在，否则任务失败）
6. 关键异常的影像截图（直接展示给用户查看）

**✅ 完成后必须询问用户**：

> "本次阅片已完成，共发现 **Y 处异常**（详见上方汇总）。**是否需要我对其中某个异常部分进行重点复核或详细说明？** 例如：指定某个结节的层面范围、对某处可疑发现进行二次确认、或针对某个 Lung-RADS 分级给出更详细的解读。"

- 如果用户回复"是"或指定了某个异常区域，则**立即对该区域重新检视对应层面的 PNG 图片**，给出更详细的分析结论，并说明是否需要更新报告。
- 如果用户回复"否"或无需复核，则任务结束。
- **此询问不属于"确认"禁令范畴**——它是在任务完成后主动提供的增值服务，不影响执行流程。

## 输入参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| input_path / --input | string | 是 | - | DICOM 文件路径或 ZIP 压缩包路径 |
| --output-dir / --output | string | 否 | 输入文件同级目录 | PNG 图片和报告的输出目录 |
| --enhance | boolean | 否 | false | 启用 Real-ESRGAN 超分辨率增强 |
| --enhance-scale | integer | 否 | 2 | 超分增强放大倍数（2 或 4） |
| --window | string | 否 | lung | 窗口类型：`lung`/`mediastinum`/`bone`/`soft_tissue`/`ggo`/`all` |
| --separate-window-dirs | boolean | 否 | true | 不同窗口类型 PNG 是否分子目录存放 |
| --mip | boolean | 否 | false | 启用 MIP（最大密度投影）重建，提高微小肺结节检出率 |
| --mip-slabs | integer | 否 | 5 | MIP slab 厚度（层数），范围 2-20 |
| --imaging-type | string | 否 | 自动检测 | 手动指定影像类型：`chest_ct`/`abdomen_ct`/`brain_mri`/`abdomen_mri`/`generic`。**转换功能支持所有类型；AI 阅片当前仅 `chest_ct` 有效** |
| --model-name | string | 否 | 无 | 阅片大模型名称，记录到 PDF 报告。宿主 AI 应自动传入自身模型名称 |
| --review-results-json | string | 否 | 无 | 已完成逐张阅片后的 `review_results.json` 路径；提供后将直接加载正式结果生成报告 |
| --strict-review | boolean | 否 | false | 启用后，若仍存在"无法识别/待检视"条目，则拒绝生成最终报告并退出 |
| --auto-review-model | string | 否 | 无 | 外部视觉模型名称；提供后会自动调用 OpenAI 兼容多模态接口逐批回填 `review_batch_templates` |
| --auto-review-api-base | string | 否 | `https://api.openai.com/v1` | 外部视觉模型的 OpenAI 兼容接口基地址 |
| --auto-review-api-key / --auto-review-api-key-env | string | 否 | 环境变量 `OPENAI_API_KEY` | 外部视觉模型 API Key，可显式传入或通过环境变量读取 |
| --auto-review-detail | string | 否 | `high` | 自动阅片时传给外部视觉模型的图片细节级别：`low` / `high` / `auto` |
| --auto-review-timeout | integer | 否 | 180 | 自动阅片单条请求超时秒数 |

## 输出

- **PNG 图片**：默认肺窗 + 自动输出纵隔窗和 GGO 窗。`--window all` 输出全部 5 种窗口
- **MIP 重建图像**（`--mip`）：存放在 `mip/` 子目录
- **PDF 报告**：医院风格 AI 辅助影像检查报告（含检查信息表、检查所见、异常影像展示、诊断意见、免责声明）
- **Markdown 报告**：与 PDF 同名同目录，方便版本控制

分目录模式（默认）下的输出结构：
```
<output_dir>/<时间戳>/
├── png/
│   ├── lung/           # 肺窗 — 检视实性结节、肺纹理
│   ├── mediastinum/    # 纵隔窗 — 验证结节密度、检查淋巴结
│   ├── ggo/            # GGO 专用窗 — ⚠️ 必须优先检视！磨玻璃结节可能仅此窗可见
│   ├── bone/           # 骨窗（--window all）
│   ├── soft_tissue/    # 软组织窗（--window all）
│   └── mip/            # MIP（--mip）
├── review_manifest.json
├── review_results_stub.json
├── review_batch_templates/
├── review_batch_filled/       # 使用 auto_review_batches.py 或 --auto-review-model 后生成
├── review_results.json        # 合并后的正式结果总表
├── dicom_report_<时间戳>.pdf   # 医院风格 PDF 报告
└── dicom_report_<时间戳>.md    # Markdown 报告
```

## 工作流程

```
输入（DICOM/ZIP）→ 影像类型识别 → 转换 PNG → [可选]超分增强 → [可选]MIP重建 → AI 全量阅片 → PDF+MD 报告
```

1. **DICOM 转 PNG**：自动检测后端（DCMTK → SimpleITK → dicom2jpg），CT 默认肺窗（WC=-600, WW=1500）+ 自动输出纵隔窗和 GGO 窗（WC=-500, WW=800）。低分辨率图像自动 Lanczos 放大到至少 1024×1024
2. **超分增强**（可选）：Real-ESRGAN 提升清晰度
3. **MIP 重建**（可选）：连续多层最大密度投影，提高 2-6mm 肺结节检出率
4. **AI 全量阅片**：逐批（15-20 张）检视全部切片，检视顺序 GGO窗→肺窗→纵隔窗，每张标注层面序号（如"第285/832层"）
5. **报告生成**：同时输出 PDF 和 Markdown 报告（格式参考 `references/AI_chest_CT_report_template.pdf`）

## 脚本说明

### 命令行调用

```bash
# 推荐写法（必须传入 --output 参数！）
python3 scripts/main.py --input <input_path> --output <output_dir> [options]

# 常用组合示例
python3 scripts/main.py --input /path/to/chest.zip --output /path/to/dicom_output_20260319_151835
python3 scripts/main.py --input /path/to/chest.zip --output /path/to/output --enhance --enhance-scale 4
python3 scripts/main.py --input /path/to/chest.zip --output /path/to/output --mip --window all
python3 scripts/main.py --input /path/to/chest.zip --output /path/to/output --model-name claude-4.6-opus --strict-review
# 注意：当前版本仅支持胸部CT，其他影像类型（abdomen_ct/brain_mri/abdomen_mri）暂不支持
# 严格模式 + 外部视觉模型自动逐批回填（OpenAI 兼容接口）
python3 scripts/main.py --input /path/to/chest.zip --output /path/to/output --strict-review --auto-review-model gpt-4.1 --auto-review-api-base https://api.openai.com/v1 --auto-review-api-key "$OPENAI_API_KEY"
# 已有接力包时，单独调用外部视觉模型自动补跑全部批次
python3 scripts/auto_review_batches.py --manifest /path/to/output/<时间戳>/review_manifest.json --model gpt-4.1 --api-base https://api.openai.com/v1 --api-key "$OPENAI_API_KEY"
# 手工/半自动模式下：每完成一批 batch_XXX.json 回填，就把结果并入总表
python3 scripts/apply_review_batch.py --manifest /path/to/output/<时间戳>/review_manifest.json --results /path/to/output/<时间戳>/review_results_stub.json --batch-json /path/to/output/<时间戳>/review_batch_templates/batch_001.json --output /path/to/output/<时间戳>/review_results_working.json
# 全部批次回填完成后，生成正式报告
python3 scripts/generate_report.py --results /path/to/output/<时间戳>/review_results.json --manifest /path/to/output/<时间戳>/review_manifest.json --output /path/to/output/<时间戳> --input-path /path/to/chest.zip --imaging-type chest_ct --model-name claude-4.6-opus
# 如确实只想生成草稿报告（例如调试），可显式允许不完整结果
python3 scripts/generate_report.py --results /path/to/output/<时间戳>/review_results_stub.json --output /path/to/output/<时间戳>/draft --input-path /path/to/chest.zip --imaging-type chest_ct --model-name claude-4.6-opus --allow-incomplete
```

### 脚本列表

| 脚本 | 说明 |
|------|------|
| `scripts/run.py` | **跨平台启动器（推荐入口）**：纯 Python，自动生成时间戳、创建输出目录、查找 Python 解释器、检测 pip 镜像，macOS/Linux/Windows 通用 |
| `scripts/pip_utils.py` | **pip 镜像感知安装模块**：自动检测网络环境，国内用户自动切换清华/阿里/华为/中科大镜像，所有 pip install 操作统一走此模块 |
| `scripts/main.py` | 主入口，编排流水线（含影像类型识别） |
| `scripts/modality_detector.py` | 影像类型自动识别与 ImagingProfile 策略模型 |
| `scripts/converter.py` | DICOM → PNG 转换（根据影像类型选择窗位策略） |
| `scripts/enhancer.py` | Real-ESRGAN 超分辨率增强 |
| `scripts/reviewer.py` | AI 阅片检视（根据影像类型选择 Prompt 模板和分级系统，并导出批次接力包） |
| `scripts/report_generator.py` | PDF 报告生成（根据影像类型动态生成报告格式） |
| `scripts/auto_review_batches.py` | **外部视觉模型自动回填工具**：逐批读取 `review_batch_templates/batch_XXX.json`，调用 OpenAI 兼容多模态接口并持续合并生成 `review_results.json` |
| `scripts/apply_review_batch.py` | **批次结果合并工具**：把 `review_batch_templates/batch_XXX.json` 中回填好的结果并入总表 |
| `scripts/generate_report.py` | **独立报告生成入口**：从阅片结果 JSON 生成 PDF（默认正式模式，会校验 manifest 并拒绝待检视条目） |
| `scripts/full_auto_review.py` | **v2.6.0 新增 - 宿主 AI 全自动化阅片**：利用宿主 AI 的多模态能力逐张检视全部切片，实现真正的零遗漏全量阅片 |
| `scripts/prompt_templates/` | 各影像类型的 Prompt 模板目录 |

## 全量阅片要点（宿主 AI 必读）

### 为什么不能抽样

一个 3mm 肺结节在 800 层 CT（1mm 层距）中仅出现在约 3 层上。如果每隔 40 层看一张（共 20 张），漏诊概率超过 **96%**。这不是效率问题，是患者安全问题。

> **绝对禁令**：无论出于何种原因（Token 限制、时间限制、上下文窗口限制），都**不得以抽样代替全量阅片**。如果模型能力不足以一次完成，则必须分多个会话完成，而不是降级为抽样。

### 高效全量阅片策略

1. **分批执行**：每批 15-20 张原始分辨率切片
2. **三窗联合**：GGO窗（优先）→ 肺窗 → 纵隔窗
3. **双侧强制扫查**：每张图片检查左右两侧
4. **反保守偏见**（见铁律第 10 条）：≥2mm 的圆形高密度影优先报告为疑似结节，禁止轻易归为"血管断面"。可疑/不确定的发现**必须标"异常"**（见铁律第 9 条）
5. **扫及区域不遗漏**：甲状腺、肝脏上段、肾上腺

### 宿主 AI 阅片产物要求

阅片完成后，宿主 AI 应在输出目录中生成：
- 结构化的阅片发现 Markdown 文件（汇总所有异常）
- 关键异常影像截图展示
- 所有产物必须在 `<output_dir>/<时间戳>/` 目录内，禁止散落到工作区根目录

## 内置资源

| 资源 | 说明 |
|------|------|
| `fonts/NotoSansSC-Regular.ttf` | Google Noto Sans SC 中文字体（SIL OFL），PDF 中文渲染 |
| `references/AI_chest_CT_report_template.pdf` | 胸部 CT 报告模板参考（**PDF 报告必须参照此格式**） |
| `references/review_strategy.md` | 详细的全量阅片策略指南（三阶段法） |

## 自修复能力

运行时自动处理环境问题，无需手动干预：

| 场景 | 自修复行为 |
|------|-----------| 
| Python 依赖缺失 | 自动 pip 安装 |
| DICOM 转换后端不可用 | 自动安装 SimpleITK 或 dicom2jpg |
| Real-ESRGAN 不可用 | 自动安装，失败则降级使用原始图片 |
| reportlab 不可用 | 自动安装，失败则降级输出纯文本报告 |
| 输入路径不存在 | 搜索同目录下相似文件并提示 |

## 环境要求

- Python 3.8+
- DICOM 转换后端（至少一种）：DCMTK（推荐）/ SimpleITK / dicom2jpg
- Real-ESRGAN（可选，超分增强）

```bash
pip install -r requirements.txt
```

## JSON Schema 统一字段（v2.4.0+）

从 v2.4.0 起，所有影像类型的 AI 阅片 JSON 输出采用统一字段：

| 字段 | 说明 | 兼容性 |
|------|------|--------|
| `classification_system` | 分级系统名称（Lung-RADS / LI-RADS / Fazekas 等） | 新增 |
| `classification_value` | 分级值（如 "2类"、"LR-3"、"2级"） | 新增 |
| `bounding_boxes` | 异常区域归一化坐标 `[{x, y, width, height}]` | 全模态通用 |
| `lung_rads` | Lung-RADS 分类值（胸部CT专用，向后兼容） | 保留 |

> 胸部CT 的 `lung_rads` 字段保留向后兼容；其他影像类型统一使用 `classification_system` + `classification_value`。

## Changelog

### v2.5.0 — CAD v2.8 假阳性深度优化（Case2 医院报告精确对照）

基于 Case2（岚天 2026-01-10 胸部CT, 1.25mm）医院报告 vs CAD 的精确对照分析——医院确认4个结节（1个GGO+3个实性），CAD输出35个候选中31个假阳性（FP率89%）。深度分析假阳性模式后实施以下优化：

**CAD v2.8 核心改动**

- **[新] 肺外缘空间惩罚**：cx距图像中心>180px→0.4, >150px→0.6, >120px→0.8。Case2有5个外缘假阳性(S#3,S#4,S#8,S#11,S#15)全在cx>370区域。保护：4个真结节cx距中心均<100px，不触发。
- **[新] 外缘+微小灶联合惩罚**：HU 80-200 + z_slices=2 + d<2.5mm + cx_dist>100 → 0.5。这是Case2中最典型的血管截面假阳性组合模式。
- **[改] GGO HU评分精细化**：-450~-380→1.0, -480~-450→0.8, -500~-480→0.5（原-500~-400全给1.0）。以-480为新分水岭。Case2真GGO HU=-409仍在最佳区间，假阳性(HU=-476~-492)显著降分。
- **[改] GGO合并距离**：15mm→20mm。Case2有3个GGO假阳性与真GGO在15mm内残留未合并。
- **[改] GGO血管排除收紧**：mean_hu惩罚阈值从-530/-510收紧到-510/-500/-480，配合HU评分形成双重降权。

**校准数据更新**
- Case2 医院报告精确匹配：GGO#21→GT1, Solid#6→GT2, Solid#17→GT3, Solid#5→GT4
- DICOM方向确认：心脏在图像左侧, cx<256=LEFT, cx>256=RIGHT

**预期效果**
- 假阳性从~31个降至~14个（FP率89%→~60%）
- 4个真结节评分保持≥0.90（0漏检）

### v2.4.8 — CAD vs AI 交叉验证（结果复核机制）

新增自动化的 **CAD vs AI 交叉验证**步骤，在 AI 阅片完成后、PDF 报告生成前自动执行：

**核心逻辑**
- 对所有 CAD 评分 ≥ 0.80 的高分候选，检查对应 z-range（±2 层容差）内的 AI 阅片结论
- 三类告警检测：
  - `missed_candidate`：CAD 高分但 AI 全标正常且无排除理由 → 可能遗漏
  - `described_but_not_reported`：AI 描述中提及可疑征象但结论为正常 → 矛盾
  - `no_coverage`：候选区域内无 AI 实际审阅结果 → 盲区
- 支持智能排除：如 AI 在描述中明确提及"血管/伪影/支气管"等排除理由，视为已审查并排除

**集成位置**
- `main.py`：阶段 3.5（阶段 3 校验之后、阶段 4 报告生成之前）
- `generate_report.py`：自动从 `cad_candidates.json` 加载 CAD 结果并执行交叉验证
- 告警结果保存为 `cross_validation.json`，并在 PDF 报告中以醒目的告警表格展示

**新增函数**
- `reviewer.cross_validate_cad_vs_review()` — 交叉验证核心逻辑
- `report_generator._build_cross_validation_section()` — PDF 报告告警章节

### v2.4.7 — 上下文窗口优化（分层 CAD hint + 精简版 prompt）

解决宿主 AI 因上下文窗口耗尽（1079 层 × 4 窗位 × ~14KB/层 prompt）导致阅片任务中断的问题：

**分层 CAD hint 注入（替代全局注入）**
- 旧方式：每层都注入全部 25 个 CAD 候选（~2.5KB），800 层 = 2MB 冗余
- 新方式：每层只注入 z-index ±5 层范围内的 CAD 候选（通常 0-3 个，~200B）
- 新增 `cad_detector.format_candidates_for_slice()` 和 `get_cad_focus_slices()` 函数

**双层 prompt 策略**
- **📌重点层**（CAD 候选 ±5 层）：使用完整版 REVIEW_PROMPT（~5000 字）+ 分层 CAD hint
- **⚡快扫层**（远离 CAD 候选）：使用新增的 REVIEW_PROMPT_LITE（~600 字），仍需逐张检视
- 精简版保留核心要素：GGO 窗优先、反保守偏见、JSON 输出格式、扫及区域
- 典型 800 层 CT：重点层 ~60-100 层，快扫层 ~700+ 层，prompt token 总量节省 ~50-70%

**结构化 CAD 数据传递**
- `main.py` 新增 `cad_candidates` 参数传递给 `reviewer.review()`
- `reviewer.py` 的 `review()` 方法新增 `cad_candidates` 参数，自动计算重点层集合
- `review_manifest.json` 中每层新增 `is_focus_layer` 和 `layer_type` 字段

**review_manifest.json 体积大幅减小**
- 旧方式：832 层 × ~14KB/层 = 11.6MB
- 新方式：100 重点层 × ~7KB + 732 快扫层 × ~1.5KB = ~1.8MB（减少 ~85%）

### v2.6.0 — 宿主 AI 全自动化全量阅片（重大更新）

**新增核心功能**

- **新增 `full_auto_review.py` 模块**：利用宿主 AI 的多模态视觉能力，实现真正的零遗漏全量阅片
  - 自动生成标准化阅片请求列表（1079层 → 1079个独立请求）
  - 支持宿主 AI 逐张分析图片并回填结果
  - 每10张自动保存进度，防止中断丢失
  - 生成标准化的 `review_results_auto.json` 结果文件
  
- **版本号升级**：2.5.0 → 2.6.0
  - `version.py` 更新版本号
  - SKILL.md 文档同步更新
  - 脚本列表新增 `full_auto_review.py` 说明

**技术改进**

- 解决了手动逐批阅片耗时过长的问题（72批次 × 15张 = 1079张）
- 保持了 dicom-doctor 的铁律要求：逐层确认、零遗漏、全量阅片
- 兼容现有 `generate_report.py` 报告生成流程

**使用方式**

```bash
# 生成全量阅片请求
python3 scripts/full_auto_review.py --manifest <path>/review_manifest.json --output <output_dir>

# 宿主 AI 逐张处理请求后，生成报告
python3 scripts/generate_report.py --results <output_dir>/review_results_auto.json --manifest <path>/review_manifest.json ...
```

### v2.4.6 — 结节聚合正则修复（兼容 第N/M层 格式）

- 修复 `_normalize_location()` 正则表达式：旧正则 `[\d\-]+` 只匹配数字和连字符，无法处理 `（第344/832层）` 这种含斜杠的层面编号格式
- 新正则支持所有常见层面编号格式：
  - `(第218-221层)` — 连字符范围
  - `（第344/832层）` — 斜杠分隔（当前层/总层数）
  - `第218层` — 无括号版本
  - `，第197-199层` — 逗号连接版本
  - `(第218,219,220层)` — 逗号列举版本
- 新增尾部残留标点清理逻辑
- 修复了"不同智能体生成的 location 字段格式不一致导致结节聚合失败"的问题（如 84 个异常层面被当成 84 个独立结节）

### v2.4.5 — PDF 报告页数控制（最多 50 页）

- 新增 `MAX_PDF_PAGES = 50` 全局常量和 `max_pages` 参数
- 异常影像展示区域增加**页数预算控制**：
  - 固定区域（标题/信息/所见/诊断/分级/免责）预估 8 页
  - 剩余页数按每张异常影像 ~1.5 页估算，计算最大可展示图片数
  - **优先展示**每个聚合结节的代表层面（置信度最高的那张，带 ★ 标记）
  - 预算有余则按顺序补充非代表层面
  - 被省略的层面用 8pt 简明汇总表替代（含 DICOM 文件名/位置/描述/置信度）
  - 汇总表末尾提示"完整异常影像请查阅输出目录下的 PNG 文件"
- `generate_report.py` 新增 `--max-pages` 命令行参数（默认 50）
- `main.py` 无需修改（`max_pages` 有默认值，零侵入兼容）

### v2.4.4 — CAD v2.7 假阳性深度优化

基于 2023-06-02 胸部CT（Case 1）vs 医院报告的复核分析——AI 报了 13 个结节，医院只报了 1 个（右肺下叶前基底段 2mm 炎性肉芽肿）。深度分析 24 个实性 + 15 个 GGO 候选的假阳性模式后：

- **实性 HU 评分进一步收紧**：mean\_hu < 20 → 0.15（Case1 大量 mean\_hu=1\~11 的假阳性，真结节最低 49）
- **新增血管模式惩罚**：(max\_hu - mean\_hu) > 200 且 mean\_hu < 50 → vessel\_penalty × 0.2\~0.4（血管横截面典型模式：中心高密度+边缘低密度被平均拉低）
- **GGO HU 下限收紧**：mean\_hu < -520 → vessel\_penalty × 0.3\~0.5（接近空气-1000，真 GGO 的 HU 一般 > -500）
- **肺尖/肺底区域惩罚**：z < 5% 或 z > 95% → vessel\_penalty × 0.4（部分容积效应/胸廓入口伪影高发区）
- **GGO 合并距离扩大**：从 8mm → 15mm（Case1 有 5 个 GGO 假阳性在 20mm 内聚集，应合并为 1 个）
- **回归验证**：5 个真结节评分全部 ≥ 0.90（Case1 GT: 0.960, Case2 GT1-4: 0.907\~0.971），假阳性肺尖 Solid#1 从 0.925 → 0.498

### v2.4.3 — 可疑发现必报铁律 + 反保守偏见铁律化

基于真实漏诊案例（2023-06-02 胸部CT，CAD score=0.96 的实性结节被 AI 描述为"可疑微小结节"但 conclusion 标为"正常"，导致未进入报告——而该结节正是医院确诊的右肺下叶前基底段 2mm×2mm 炎性肉芽肿）：

- **铁律新增第 9 条**：可疑/疑似/不确定的发现**必须标记为"异常"**并进入报告。严禁"描述了但不报"——这是最危险的漏诊模式
- **铁律新增第 10 条**：反保守偏见从"阅片策略建议"提升为铁律。≥2mm 圆形高密度影必须优先报告为疑似结节，附具体判定规则
- **铁律第 7 条增强**：CAD 候选验证从二元判定（结节/血管）改为**三元判定**（确认结节✅ / 可疑待定⚠️ / 明确排除❌），可疑待定必须标异常进报告
- 核心原则：**宁可多报一个可疑发现让医生复核，也不能吞掉一个可能的结节**

### v2.4.2 — 报告两维度数据完整性铁律

- **铁律新增第 8 条**：报告必须同时包含结节维度（聚合汇总）和层面维度（逐层明细），缺一不可
- 第 4 步汇总说明详列两维度各自必须包含的字段清单
- 第 6 步报告格式明确标注两维度"缺一不可"

### v2.4.1 — CAD v2.6 假阳性评分优化

基于 2 例 5 GT 复核校准数据（含 30 实性 + 15 GGO 候选的假阳性模式分析）：

- **实性 HU 评分收紧**：80\~200→1.0, 50\~80→0.6, 20\~50→0.4, -50\~20→0.3（原 20\~200 全给 1.0）
- **GGO HU 评分收紧**：-500\~-400→1.0, -550\~-500→0.7, -600\~-550→0.5（原 -600\~-400 全给 1.0）
- **GGO elongation 惩罚**：>1.8→0.3, >1.5 且 voxels<30→0.5（真 GGO elong=1.16, 假阳性 elong=1.4\~2.0+）
- 空间聚簇惩罚暂缓（统计基础不足）
- 回归验证 4 个真结节全部通过

### v2.4.0 — 多影像类型架构统一（Phase 0）

**统一 JSON Schema**
- 所有影像类型 prompt 模板统一输出 `classification_system` + `classification_value` + `bounding_boxes` 字段
- 胸部CT 保留 `lung_rads` 字段向后兼容
- `reviewer.py` 解析逻辑支持新旧字段名双向回退

**反保守偏见铁律全模态推广**
- 腹部CT：肝脏低密度灶反保守偏见（Couinaud 8段扫查 + 易漏诊区域清单）
- 头颅MRI：FLAIR高信号灶反保守偏见（DWI真假弥散受限鉴别 + 基底节/丘脑/脑桥腔梗高发区）
- 腹部MRI：T2高信号灶+DWI弥散受限反保守偏见（多序列组合判断表）
- 通用模式：通用反保守偏见铁律

**报告分级参考表扩展**
- 新增 LI-RADS 参考表（LR-1 到 LR-TIV，含完整处理建议）
- 新增 Fazekas 白质病变分级参考表（0-3级）
- 报告生成器根据 `ImagingProfile.classification_system` 自动选择展示哪种参考表

**诊断意见随访建议**
- LI-RADS：LR-1 无需随访 → LR-5 立即治疗，7档具体建议
- Fazekas：0级正常 → 3级重度+神经内科随访，4档具体建议
- PDF 和 Markdown 报告同步支持

**Prompt 模板深度大幅提升**
- 腹部CT：4.25KB → ~14KB（+Couinaud扫查/LI-RADS完整表/易漏诊区域/CAD占位符）
- 头颅MRI：3.83KB → ~12KB（+Fazekas分级/MRI信号分析/解剖分区扩展/CAD占位符）
- 腹部MRI：3.54KB → ~10KB（+多序列组合判断表/LI-RADS完整表/CAD占位符）
- 通用模式：2.2KB → ~4KB（+反保守偏见铁律/bounding_boxes/CAD占位符）

**头颅MRI ImagingProfile 增强**
- `classification_system` 从空改为 `"Fazekas"`
- `report_sections` 新增"基底节及丘脑"和"扫及区域"

**CAD 预留接口**
- 所有非胸部CT模板新增 `{cad_hint}` 占位符，为未来 CAD 候选注入预留

## ⚠️ 免责声明

本 skill 生成的报告由 AI 辅助生成，**仅供参考，不构成医学诊断**。如有疑问，请及时咨询专业医生。
