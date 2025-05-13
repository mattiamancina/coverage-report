#!/usr/bin/env php
<?php
/**
 * parse_cobertura.php
 * Parses a Cobertura XML coverage report and outputs a Markdown summary.
 */
if ($argc !== 3) {
    fwrite(STDERR, "Usage: parse_cobertura.php <coverage.xml> <output.md>\n");
    exit(1);
}
$xmlPath = $argv[1];
$outputPath = $argv[2];
if (!file_exists($xmlPath)) {
    fwrite(STDERR, "Error: XML report file not found at '$xmlPath'\n");
    exit(1);
}
$xml = simplexml_load_file($xmlPath);
if ($xml === false) {
    fwrite(STDERR, "Error: Failed to parse XML file\n");
    exit(1);
}
// Overall coverage
$lineRate = ((float)($xml['line-rate'] ?? 0)) * 100;
$branchRate = ((float)($xml['branch-rate'] ?? 0)) * 100;
// Per-file classes coverage
$classes = [];
foreach ($xml->packages->package as $package) {
    foreach ($package->classes->class as $class) {
        $filename = (string)$class['filename'];
        $lr = ((float)($class['line-rate'] ?? 0)) * 100;
        $br = ((float)($class['branch-rate'] ?? 0)) * 100;
        $classes[] = ['filename' => $filename, 'line' => $lr, 'branch' => $br];
    }
}
// Sort classes by filename
usort($classes, function($a, $b) { return strcmp($a['filename'], $b['filename']); });
// Generate Markdown
$md = [];
$md[] = '# Coverage Report';
$md[] = '';
$md[] = "**Overall Line Coverage:** " . number_format($lineRate, 2) . "%  ";
$md[] = "**Overall Branch Coverage:** " . number_format($branchRate, 2) . "%  ";
$md[] = '';
$md[] = '## File Coverage  ';
$md[] = '';
$md[] = '| File | Line Coverage | Branch Coverage |';
$md[] = '| ---- | ------------- | --------------- |';
foreach ($classes as $c) {
    $md[] = sprintf(
        '| %s | %.2f%% | %.2f%% |',
        $c['filename'],
        $c['line'],
        $c['branch']
    );
}
// Write to output file
file_put_contents($outputPath, implode("\n", $md) . "\n");
exit(0);
?>