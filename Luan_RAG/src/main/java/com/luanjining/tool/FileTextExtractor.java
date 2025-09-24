package com.luanjining.tool;

import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.text.PDFTextStripper;
import org.apache.poi.xwpf.usermodel.XWPFDocument;
import org.apache.poi.xwpf.extractor.XWPFWordExtractor;
import org.apache.poi.hwpf.HWPFDocument;
import org.apache.poi.hwpf.extractor.WordExtractor;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;
import org.apache.poi.hssf.usermodel.HSSFWorkbook;
import org.apache.poi.ss.usermodel.*;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;

/**
 * 文件文本提取器
 * 支持多种文件格式的文本提取
 */
public class FileTextExtractor {

    /**
     * 根据文件类型提取文本内容
     * @param file 要提取文本的文件
     * @return 提取出的文本内容
     * @throws Exception 提取失败时抛出异常
     */
    public String extractText(File file) throws Exception {
        if (file == null || !file.exists()) {
            throw new FileNotFoundException("文件不存在");
        }

        String fileName = file.getName().toLowerCase();

        if (fileName.endsWith(".txt") || fileName.endsWith(".md")) {
            return extractFromTextFile(file);
        } else if (fileName.endsWith(".pdf")) {
            return extractFromPdf(file);
        } else if (fileName.endsWith(".docx")) {
            return extractFromDocx(file);
        } else if (fileName.endsWith(".doc")) {
            return extractFromDoc(file);
        } else if (fileName.endsWith(".xlsx")) {
            return extractFromXlsx(file);
        } else if (fileName.endsWith(".xls")) {
            return extractFromXls(file);
        } else if (fileName.endsWith(".csv")) {
            return extractFromCsv(file);
        } else {
            // 默认当作文本文件处理
            return extractFromTextFile(file);
        }
    }

    /**
     * 提取纯文本文件
     */
    private String extractFromTextFile(File file) throws IOException {
        return Files.readString(file.toPath(), StandardCharsets.UTF_8);
    }

    /**
     * 提取PDF文本
     */
    private String extractFromPdf(File file) throws IOException {
        try (PDDocument document = PDDocument.load(file)) {
            PDFTextStripper stripper = new PDFTextStripper();
            return stripper.getText(document);
        }
    }

    /**
     * 提取DOCX文本
     */
    private String extractFromDocx(File file) throws IOException {
        try (FileInputStream fis = new FileInputStream(file);
             XWPFDocument document = new XWPFDocument(fis);
             XWPFWordExtractor extractor = new XWPFWordExtractor(document)) {
            return extractor.getText();
        }
    }

    /**
     * 提取DOC文本
     */
    private String extractFromDoc(File file) throws IOException {
        try (FileInputStream fis = new FileInputStream(file);
             HWPFDocument document = new HWPFDocument(fis);
             WordExtractor extractor = new WordExtractor(document)) {
            return extractor.getText();
        }
    }

    /**
     * 提取XLSX文本
     */
    private String extractFromXlsx(File file) throws IOException {
        StringBuilder text = new StringBuilder();
        try (FileInputStream fis = new FileInputStream(file);
             XSSFWorkbook workbook = new XSSFWorkbook(fis)) {

            for (Sheet sheet : workbook) {
                text.append("Sheet: ").append(sheet.getSheetName()).append("\n");
                for (Row row : sheet) {
                    for (Cell cell : row) {
                        String cellValue = getCellValueAsString(cell);
                        if (!cellValue.trim().isEmpty()) {
                            text.append(cellValue).append("\t");
                        }
                    }
                    text.append("\n");
                }
                text.append("\n");
            }
        }
        return text.toString();
    }

    /**
     * 提取XLS文本
     */
    private String extractFromXls(File file) throws IOException {
        StringBuilder text = new StringBuilder();
        try (FileInputStream fis = new FileInputStream(file);
             HSSFWorkbook workbook = new HSSFWorkbook(fis)) {

            for (Sheet sheet : workbook) {
                text.append("Sheet: ").append(sheet.getSheetName()).append("\n");
                for (Row row : sheet) {
                    for (Cell cell : row) {
                        String cellValue = getCellValueAsString(cell);
                        if (!cellValue.trim().isEmpty()) {
                            text.append(cellValue).append("\t");
                        }
                    }
                    text.append("\n");
                }
                text.append("\n");
            }
        }
        return text.toString();
    }

    /**
     * 提取CSV文本
     */
    private String extractFromCsv(File file) throws IOException {
        return Files.readString(file.toPath(), StandardCharsets.UTF_8);
    }

    /**
     * 获取Excel单元格的字符串值
     */
    private String getCellValueAsString(Cell cell) {
        if (cell == null) return "";

        switch (cell.getCellType()) {
            case STRING:
                return cell.getStringCellValue();
            case NUMERIC:
                if (DateUtil.isCellDateFormatted(cell)) {
                    return cell.getDateCellValue().toString();
                } else {
                    // 格式化数字，避免科学计数法
                    double numericValue = cell.getNumericCellValue();
                    if (numericValue == (long) numericValue) {
                        return String.valueOf((long) numericValue);
                    } else {
                        return String.valueOf(numericValue);
                    }
                }
            case BOOLEAN:
                return String.valueOf(cell.getBooleanCellValue());
            case FORMULA:
                return cell.getCellFormula();
            case BLANK:
                return "";
            default:
                return "";
        }
    }

    /**
     * 检查文件是否为支持的类型
     */
    public boolean isSupportedFileType(File file) {
        if (file == null || !file.exists()) {
            return false;
        }

        String fileName = file.getName().toLowerCase();
        return fileName.endsWith(".txt") ||
                fileName.endsWith(".md") ||
                fileName.endsWith(".pdf") ||
                fileName.endsWith(".docx") ||
                fileName.endsWith(".doc") ||
                fileName.endsWith(".xlsx") ||
                fileName.endsWith(".xls") ||
                fileName.endsWith(".csv");
    }

    /**
     * 提取文件中的文本内容
     * @param file 要提取文本的文件
     * @return 提取出的文本内容
     * @throws Exception 提取失败时抛出异常
     */
    public String getExtractedText(File file) throws Exception {
        FileTextExtractor textExtractor = new FileTextExtractor();

        // 1. 检查文件是否存在
        if (file == null || !file.exists()) {
            throw new IllegalArgumentException("文件不存在: " + (file != null ? file.getPath() : "null"));
        }

        // 2. 检查是否为支持的文件类型
        if (!textExtractor.isSupportedFileType(file)) {
            throw new UnsupportedOperationException("不支持的文件类型: " + file.getName());
        }

        // 3. 提取文本
        String extractedText = textExtractor.extractText(file);

        // 4. 检查提取的文本是否为空
        if (extractedText == null || extractedText.trim().isEmpty()) {
            throw new IllegalArgumentException("文件内容为空或无法提取文本");
        }

        // 5. 清理文本，去除多余的空白字符
        extractedText = extractedText.replaceAll("[\\s\\r\\n\\t]+", "");
        return extractedText;
    }


}