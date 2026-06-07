param(
    [string]$ExistingPath = ".\data\vocabulary\vocabulary_words.json",
    [string]$FreeDictTeiPath = ".\freedict-eng-rus.tei",
    [string]$OutputPath = ".\data\vocabulary\vocabulary_words.json"
)

$ErrorActionPreference = "Stop"

function New-VocabId([string]$english) {
    $id = $english.Trim().ToLowerInvariant()
    $id = $id.Replace(" ", "-").Replace("'", "").Replace("/", "-").Replace(".", "").Replace(",", "")
    return $id
}

function Normalize-Text([string]$value) {
    if ($null -eq $value) {
        return ""
    }
    return [regex]::Replace($value.Trim(), "\s+", " ")
}

function Is-LearnableEnglish([string]$value) {
    if ([string]::IsNullOrWhiteSpace($value)) {
        return $false
    }
    $value = $value.Trim()
    if ($value.Length -gt 80) {
        return $false
    }
    if ($value -notmatch "[A-Za-z]") {
        return $false
    }
    if ($value -match "[<>{}\[\]|]") {
        return $false
    }
    return $true
}

$existing = Get-Content -LiteralPath $ExistingPath -Encoding UTF8 -Raw | ConvertFrom-Json
$seen = @{}
$output = New-Object System.Collections.Generic.List[object]

foreach ($word in $existing) {
    $id = New-VocabId $word.english
    if ($seen.ContainsKey($id)) {
        continue
    }
    $word.id = $id
    $seen[$id] = $true
    $output.Add($word)
}

$settings = [System.Xml.XmlReaderSettings]::new()
$settings.DtdProcessing = [System.Xml.DtdProcessing]::Ignore
$settings.IgnoreComments = $true
$settings.IgnoreWhitespace = $true

$reader = [System.Xml.XmlReader]::Create((Resolve-Path -LiteralPath $FreeDictTeiPath), $settings)
$added = 0

try {
    while ($reader.Read()) {
        if ($reader.NodeType -ne [System.Xml.XmlNodeType]::Element -or $reader.LocalName -ne "entry") {
            continue
        }

        $entry = $reader.ReadSubtree()
        $english = ""
        $pos = ""
        $translations = New-Object System.Collections.Generic.List[string]
        $transDepth = -1

        try {
            while ($entry.Read()) {
                if ($entry.NodeType -eq [System.Xml.XmlNodeType]::Element) {
                    if ($entry.LocalName -eq "orth" -and $english -eq "") {
                        $english = Normalize-Text $entry.ReadElementContentAsString()
                        continue
                    }
                    if ($entry.LocalName -eq "pos" -and $pos -eq "") {
                        $pos = Normalize-Text $entry.ReadElementContentAsString()
                        continue
                    }
                    if ($entry.LocalName -eq "cit" -and $entry.GetAttribute("type") -eq "trans") {
                        $transDepth = $entry.Depth
                        continue
                    }
                    if ($entry.LocalName -eq "quote" -and $transDepth -ge 0) {
                        $translation = Normalize-Text $entry.ReadElementContentAsString()
                        if ($translation -ne "" -and -not $translations.Contains($translation)) {
                            $translations.Add($translation)
                        }
                        continue
                    }
                }
                if ($entry.NodeType -eq [System.Xml.XmlNodeType]::EndElement -and $entry.LocalName -eq "cit" -and $entry.Depth -eq $transDepth) {
                    $transDepth = -1
                }
            }
        }
        finally {
            $entry.Close()
        }

        if (-not (Is-LearnableEnglish $english) -or $translations.Count -eq 0) {
            continue
        }
        $id = New-VocabId $english
        if ($id -eq "" -or $seen.ContainsKey($id)) {
            continue
        }

        $russian = ($translations | Select-Object -First 3) -join "; "
        if ($russian.Length -gt 160) {
            $russian = $russian.Substring(0, 157).TrimEnd() + "..."
        }

        $output.Add([pscustomobject]@{
            english = $english
            russian = $russian
            level = "C2"
            topic = "freedict"
            part_of_speech = $pos
            source = "FreeDict/WikDict eng-rus 2025.11.23 (CC BY-SA 3.0)"
            id = $id
        })
        $seen[$id] = $true
        $added++
    }
}
finally {
    $reader.Close()
}

$json = $output | ConvertTo-Json -Depth 5
$outputDirectory = Split-Path -Parent $OutputPath
if ($outputDirectory -and -not (Test-Path -LiteralPath $outputDirectory)) {
    New-Item -ItemType Directory -Force -Path $outputDirectory | Out-Null
}
[System.IO.File]::WriteAllText($OutputPath, $json + [Environment]::NewLine, [System.Text.UTF8Encoding]::new($false))

Write-Host "Existing kept: $($output.Count - $added)"
Write-Host "Added from FreeDict: $added"
Write-Host "Total vocabulary entries: $($output.Count)"
