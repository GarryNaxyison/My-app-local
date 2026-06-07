Add-Type -AssemblyName System.Drawing

$source = @"
using System;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.Drawing.Imaging;
using System.Drawing.Text;
using System.IO;

public static class PlanArtGenerator
{
    const int W = 1280;
    const int H = 720;

    public static void Generate(string outputDir)
    {
        Directory.CreateDirectory(outputDir);
        DrawFree(Path.Combine(outputDir, "plan-free.png"));
        DrawPremium(Path.Combine(outputDir, "plan-premium.png"));
        DrawPlatinum(Path.Combine(outputDir, "plan-platinum.png"));
    }

    static Color Hex(string value)
    {
        return ColorTranslator.FromHtml(value);
    }

    static Color A(int alpha, string value)
    {
        Color c = Hex(value);
        return Color.FromArgb(alpha, c.R, c.G, c.B);
    }

    static GraphicsPath RoundRect(RectangleF rect, float radius)
    {
        float d = radius * 2f;
        GraphicsPath path = new GraphicsPath();
        path.AddArc(rect.X, rect.Y, d, d, 180, 90);
        path.AddArc(rect.Right - d, rect.Y, d, d, 270, 90);
        path.AddArc(rect.Right - d, rect.Bottom - d, d, d, 0, 90);
        path.AddArc(rect.X, rect.Bottom - d, d, d, 90, 90);
        path.CloseFigure();
        return path;
    }

    static void FillRound(Graphics g, Brush brush, RectangleF rect, float radius)
    {
        using (GraphicsPath path = RoundRect(rect, radius)) g.FillPath(brush, path);
    }

    static void StrokeRound(Graphics g, Pen pen, RectangleF rect, float radius)
    {
        using (GraphicsPath path = RoundRect(rect, radius)) g.DrawPath(pen, path);
    }

    static void Setup(Graphics g)
    {
        g.SmoothingMode = SmoothingMode.AntiAlias;
        g.InterpolationMode = InterpolationMode.HighQualityBicubic;
        g.PixelOffsetMode = PixelOffsetMode.HighQuality;
        g.TextRenderingHint = TextRenderingHint.AntiAliasGridFit;
    }

    static void Background(Graphics g, string a, string b, string c)
    {
        using (LinearGradientBrush bg = new LinearGradientBrush(new Rectangle(0, 0, W, H), Hex("#fbfdff"), Hex("#eef6fb"), 35f))
        {
            ColorBlend blend = new ColorBlend();
            blend.Positions = new float[] { 0f, .48f, 1f };
            blend.Colors = new Color[] { Hex("#fbfdff"), Hex(b), Hex("#f4f8fc") };
            bg.InterpolationColors = blend;
            g.FillRectangle(bg, 0, 0, W, H);
        }
        using (SolidBrush brush = new SolidBrush(A(42, a))) g.FillEllipse(brush, -140, -120, 420, 360);
        using (SolidBrush brush = new SolidBrush(A(44, b))) g.FillEllipse(brush, 910, 30, 420, 340);
        using (SolidBrush brush = new SolidBrush(A(40, c))) g.FillEllipse(brush, 800, 500, 420, 280);

        using (Pen pen = new Pen(A(30, "#2457d6"), 1f))
        {
            for (int x = -80; x < W + 80; x += 64) g.DrawLine(pen, x, 0, x + 260, H);
        }
        using (Pen pen = new Pen(A(24, "#008f7a"), 1f))
        {
            for (int x = 0; x < W + 180; x += 86) g.DrawLine(pen, x, H, x - 260, 0);
        }
    }

    static void GlassPanel(Graphics g, RectangleF rect)
    {
        using (SolidBrush brush = new SolidBrush(Color.FromArgb(220, 255, 255, 255))) FillRound(g, brush, rect, 30);
        using (Pen pen = new Pen(Color.FromArgb(90, 36, 87, 214), 2f)) StrokeRound(g, pen, rect, 30);
    }

    static void DrawHeader(Graphics g, string title, string subtitle, string color)
    {
        using (Font titleFont = new Font("Segoe UI", 52, FontStyle.Bold, GraphicsUnit.Pixel))
        using (Font subFont = new Font("Segoe UI", 22, FontStyle.Bold, GraphicsUnit.Pixel))
        using (SolidBrush ink = new SolidBrush(Hex("#111827")))
        using (SolidBrush muted = new SolidBrush(A(180, "#475569")))
        using (SolidBrush accent = new SolidBrush(Hex(color)))
        {
            g.DrawString(title, titleFont, ink, 78, 68);
            g.DrawString(subtitle, subFont, muted, 82, 132);
            FillRound(g, accent, new RectangleF(82, 178, 178, 12), 6);
        }
    }

    static void DrawBadge(Graphics g, string text, string color)
    {
        RectangleF r = new RectangleF(930, 64, 238, 54);
        using (SolidBrush brush = new SolidBrush(A(235, "#ffffff"))) FillRound(g, brush, r, 27);
        using (Pen pen = new Pen(A(92, color), 2f)) StrokeRound(g, pen, r, 27);
        using (Font font = new Font("Segoe UI", 18, FontStyle.Bold, GraphicsUnit.Pixel))
        using (SolidBrush textBrush = new SolidBrush(Hex(color))) g.DrawString(text, font, textBrush, r.X + 30, r.Y + 15);
    }

    static void Save(Bitmap bmp, string path)
    {
        bmp.Save(path, ImageFormat.Png);
    }

    static void DrawFree(string path)
    {
        using (Bitmap bmp = new Bitmap(W, H))
        using (Graphics g = Graphics.FromImage(bmp))
        {
            Setup(g);
            Background(g, "#2457d6", "#e9fbf7", "#f0b348");
            DrawHeader(g, "Free", "Start learning with calm daily practice", "#008f7a");
            DrawBadge(g, "CORE ACCESS", "#2457d6");

            using (SolidBrush hill = new SolidBrush(A(70, "#008f7a")))
            {
                PointF[] points = { new PointF(0, 560), new PointF(230, 470), new PointF(520, 530), new PointF(810, 440), new PointF(W, 500), new PointF(W, H), new PointF(0, H) };
                g.FillPolygon(hill, points);
            }

            GlassPanel(g, new RectangleF(240, 230, 560, 312));
            using (Pen spine = new Pen(A(120, "#2457d6"), 5f))
            using (SolidBrush page = new SolidBrush(Color.FromArgb(246, 255, 255, 255)))
            using (SolidBrush page2 = new SolidBrush(A(238, "#eef6ff")))
            {
                FillRound(g, page, new RectangleF(324, 318, 222, 142), 12);
                FillRound(g, page2, new RectangleF(546, 318, 222, 142), 12);
                g.DrawLine(spine, 546, 324, 546, 466);
                using (Pen line = new Pen(Hex("#2457d6"), 7f))
                {
                    line.StartCap = line.EndCap = LineCap.Round;
                    g.DrawLine(line, 370, 356, 496, 356);
                    g.DrawLine(line, 370, 388, 474, 388);
                    g.DrawLine(line, 594, 356, 712, 356);
                    g.DrawLine(line, 594, 388, 690, 388);
                }
            }
            using (SolidBrush card = new SolidBrush(Hex("#008f7a"))) FillRound(g, card, new RectangleF(820, 285, 196, 132), 24);
            using (Pen white = new Pen(Color.FromArgb(220, 255, 255, 255), 9f))
            {
                white.StartCap = white.EndCap = LineCap.Round;
                g.DrawLine(white, 858, 329, 956, 329);
                g.DrawLine(white, 858, 369, 920, 369);
            }
            using (SolidBrush glow = new SolidBrush(A(176, "#f0b348"))) g.FillEllipse(glow, 950, 206, 148, 148);
            using (SolidBrush dot = new SolidBrush(Hex("#33d5b2"))) g.FillEllipse(dot, 725, 258, 42, 42);
            using (Font plus = new Font("Segoe UI", 24, FontStyle.Bold, GraphicsUnit.Pixel))
            using (SolidBrush whiteText = new SolidBrush(Color.White)) g.DrawString("+", plus, whiteText, 735, 262);
            Save(bmp, path);
        }
    }

    static void DrawPremium(string path)
    {
        using (Bitmap bmp = new Bitmap(W, H))
        using (Graphics g = Graphics.FromImage(bmp))
        {
            Setup(g);
            Background(g, "#2457d6", "#edf7ff", "#7a4fd1");
            DrawHeader(g, "Premium", "AI tutor, voice, photos and vocabulary flow", "#2457d6");
            DrawBadge(g, "VOICE + PHOTO", "#7a4fd1");

            using (Pen orbit = new Pen(A(92, "#2457d6"), 4f))
            {
                g.DrawEllipse(orbit, 304, 204, 470, 300);
                g.DrawEllipse(orbit, 344, 158, 390, 390);
            }
            GlassPanel(g, new RectangleF(348, 220, 446, 286));
            using (SolidBrush mic = new SolidBrush(Hex("#2457d6"))) FillRound(g, mic, new RectangleF(430, 292, 78, 150), 39);
            using (Pen micPen = new Pen(Hex("#2457d6"), 12f))
            {
                micPen.StartCap = micPen.EndCap = LineCap.Round;
                g.DrawLine(micPen, 469, 456, 469, 510);
                g.DrawLine(micPen, 424, 510, 514, 510);
            }
            using (SolidBrush chat = new SolidBrush(Color.FromArgb(248, 255, 255, 255))) FillRound(g, chat, new RectangleF(545, 292, 220, 138), 24);
            using (Pen line = new Pen(Hex("#008f7a"), 8f))
            {
                line.StartCap = line.EndCap = LineCap.Round;
                g.DrawLine(line, 586, 326, 710, 326);
                g.DrawLine(line, 586, 360, 682, 360);
                g.DrawLine(line, 586, 394, 720, 394);
            }
            using (SolidBrush photo = new SolidBrush(Hex("#7a4fd1"))) FillRound(g, photo, new RectangleF(792, 250, 168, 124), 24);
            using (Pen chart = new Pen(Color.White, 9f))
            {
                chart.StartCap = chart.EndCap = LineCap.Round;
                g.DrawLine(chart, 826, 320, 858, 290);
                g.DrawLine(chart, 858, 290, 890, 316);
                g.DrawLine(chart, 890, 316, 924, 278);
            }
            using (SolidBrush teal = new SolidBrush(Hex("#33d5b2"))) g.FillEllipse(teal, 342, 248, 46, 46);
            using (SolidBrush gold = new SolidBrush(Hex("#f0b348"))) g.FillEllipse(gold, 900, 452, 62, 62);
            Save(bmp, path);
        }
    }

    static void DrawPlatinum(string path)
    {
        using (Bitmap bmp = new Bitmap(W, H))
        using (Graphics g = Graphics.FromImage(bmp))
        {
            Setup(g);
            Background(g, "#1b43ad", "#fff8e8", "#008f7a");
            DrawHeader(g, "Platinum", "Maximum limits for intense language sprints", "#1b43ad");
            DrawBadge(g, "MAX LIMITS", "#c78413");

            using (SolidBrush beam = new SolidBrush(A(64, "#f0b348")))
            {
                PointF[] points = { new PointF(175, 620), new PointF(520, 210), new PointF(670, 210), new PointF(1030, 620) };
                g.FillPolygon(beam, points);
            }
            GlassPanel(g, new RectangleF(330, 224, 570, 318));
            using (SolidBrush crown = new SolidBrush(Hex("#f0b348")))
            using (Pen dark = new Pen(Hex("#1b43ad"), 10f))
            {
                PointF[] crownPts = {
                    new PointF(424, 438), new PointF(468, 276), new PointF(560, 380),
                    new PointF(642, 240), new PointF(724, 380), new PointF(816, 276),
                    new PointF(860, 438)
                };
                g.FillPolygon(crown, crownPts);
                g.DrawLines(dark, crownPts);
                FillRound(g, new SolidBrush(Hex("#1b43ad")), new RectangleF(420, 438, 444, 62), 14);
            }
            using (SolidBrush gem = new SolidBrush(Color.FromArgb(238, 255, 255, 255))) g.FillEllipse(gem, 606, 250, 72, 72);
            using (Pen plus = new Pen(Hex("#1b43ad"), 10f))
            {
                plus.StartCap = plus.EndCap = LineCap.Round;
                g.DrawLine(plus, 642, 270, 642, 302);
                g.DrawLine(plus, 626, 286, 658, 286);
            }
            using (Pen network = new Pen(A(132, "#008f7a"), 5f))
            {
                network.StartCap = network.EndCap = LineCap.Round;
                g.DrawLine(network, 282, 308, 378, 250);
                g.DrawLine(network, 900, 300, 1015, 236);
                g.DrawLine(network, 270, 500, 374, 530);
                g.DrawLine(network, 906, 510, 1030, 548);
            }
            using (SolidBrush blue = new SolidBrush(Hex("#2457d6")))
            using (SolidBrush teal = new SolidBrush(Hex("#008f7a")))
            {
                g.FillEllipse(blue, 252, 286, 36, 36);
                g.FillEllipse(teal, 1002, 220, 42, 42);
                g.FillEllipse(teal, 248, 488, 44, 44);
                g.FillEllipse(blue, 1018, 534, 34, 34);
            }
            Save(bmp, path);
        }
    }
}
"@

Add-Type -TypeDefinition $source -ReferencedAssemblies System.Drawing
[PlanArtGenerator]::Generate((Join-Path (Get-Location) "web\assets"))
