# paschen-pd

气体击穿的帕邢定律（Paschen law）计算工具。

帕邢定律描述了均匀电场中气体间隙的击穿电压随「压强 × 间距」乘积（pd）变化的规律：存在一个极小值，低于它时无论间隙多小都无法击穿。本项目用经典 Townsend 形式实现：

```
V_b(pd) = B·p·d / ( ln(A·p·d) − ln(ln(1 + 1/γ)) )
```

其中 pd 以 Torr·cm 为单位。A、B 为气体常数，γ 为二次电子发射系数。对干燥空气常用取值 A = 15 cm⁻¹·Torr⁻¹、B = 365 V·cm⁻¹·Torr⁻¹、γ = 0.01。

## 包含内容

- `internal/paschen` — 击穿电压、极小值、约化场强计算。
- `internal/sweep` — 在 pd 区间内采样帕邢曲线（支持对数间隔）。
- `internal/api` — HTTP 接口，提供网页控制台与 JSON 接口。

## 用法

### 命令行 / HTTP 服务

```bash
go run . -http :8080            # 在 :8080 提供网页控制台与 /api
go run . -http :8080 -example-dir example
```

启动后在浏览器打开 http://localhost:8080 即可交互；网页会载入内置示例（1 mm 间隙、760 Torr 干燥空气），调用 `/api/breakdown` 与 `/api/sweep` 并绘制曲线。所有计算都在本地进程内完成。

### HTTP 接口

- `POST /api/breakdown` — 计算单点击穿电压。
  请求：`{"p":760,"d":1,"A":15,"B":365,"gamma":0.01}`（p 单位 Torr，d 单位 mm，A/B/γ 省略时取空气默认值）。
  返回：`{"pd_torr_cm":76,"voltage":11445.5,"pd_min":18.3,"v_min":6680.6,"breakdownable":true, ...}`。
- `POST /api/sweep` — 在固定压力下扫描 d 区间，返回曲线采样点。
  请求：`{"p":760,"d_min":0.05,"d_max":20,"points":240,"log":true}`。
- `/example/` 提供示例 JSON 文件。

## 结论

对 760 Torr、1 mm 间隙的干燥空气，计算得击穿电压约 11.4 kV，帕邢极小值出现在 pd ≈ 18.3 Torr·cm（V ≈ 6.68 kV）。这与实验观测「低气压或极窄间隙反而更难击穿」的规律一致。

## 构建与测试

```bash
go build ./...
go test ./...
```

## 许可

见 LICENSE。
