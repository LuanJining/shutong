package com.luanjining;

import com.luanjining.tool.FileTextExtractor;
import com.luanjining.tool.QueryTool;
import kong.unirest.HttpResponse;
import kong.unirest.Unirest;

import java.io.FileNotFoundException;
import java.util.Properties;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.util.Properties;

public class API_demo {

    private static String Authorization_dataset;
    private static String Authorization_app;
    private static String user_id;
    private FileTextExtractor fileTextExtractor;

    /**
     * 构造函数，加载配置文件并初始化常量
     * @throws IOException 如果配置文件加载失败则抛出异常
     */
    public API_demo() throws IOException {
        Properties properties = new Properties();
        // 加载配置文件
        FileInputStream inputStream = new FileInputStream("message.properties");
        properties.load(inputStream);

        // 读取配置文件中的值并赋值给常量
        Authorization_dataset = properties.getProperty("Authorization_dataset");
        Authorization_app = properties.getProperty("Authorization_app");
        user_id = properties.getProperty("user_id");

        inputStream.close();
    }

    /**
     * 创建数据集
     * @param name 数据集名称
     * @param description 数据集描述
     * @return HTTP响应
     */
    public HttpResponse<String> create_dataset(String name, String description) {
        return Unirest.post("https://api.dify.ai/v1/datasets")
                .header("Authorization", Authorization_dataset)
                .header("Content-Type", "application/json")
                .body("{\n  \"name\": \"" + name + "\",\n  \"description\": \"" + description + "\"\n}")
                .asString();
    }

    /**
     * 创建文档
     * @param datasetId 数据集ID
     * @param title 文档标题
     * @param file 要上传的文件（支持纯文本、PDF、DOCX、DOC、ELSX、XLS、csv格式）
     * @return HTTP响应
     */
    public HttpResponse<String> create_document(String datasetId, String title, File file) throws Exception {

        fileTextExtractor = new FileTextExtractor();
        return Unirest.post("https://api.dify.ai/v1/datasets/"+datasetId+"/document/create-by-text")
                .header("Authorization", Authorization_dataset)
                .header("Content-Type", "application/json")
                .body("{\n  \"name\": \""+title+"\",\n  \"text\": \""+fileTextExtractor.getExtractedText(file)+"\",\n  \"indexing_technique\": \"high_quality\",\n  \"doc_form\": \"text_model\",\n  \"doc_language\": \"中文\",\n  \"process_rule\": {\n    \"mode\": \"automatic\",\n    \"rules\": {\n      \"pre_processing_rules\": [\n        {\n          \"id\": \"remove_extra_spaces\",\n          \"enabled\": true\n        }\n      ],\n      \"segmentation\": {\n        \"separator\": \"###\",\n        \"max_tokens\": 500\n      }\n    }\n  },\n  \"retrieval_model\": {\n    \"search_method\": \"hybrid_search\",\n    \"reranking_enable\": true,\n    \"top_k\": 5,\n    \"score_threshold_enabled\": true,\n    \"score_threshold\": 0.8,\n    \"weights\":{\n      \"semantic\": 0.5,\n      \"keyword\": 0.5\n    }\n  },\n  \"embedding_model\": \"text-embedding-ada-002\",\n  \"embedding_model_provider\": \"openai\"\n}")
                .asString();

    }

    /**
     * 更新文档
     * @param datasetId 数据集ID
     * @param documentId 文档ID
     * @param new_title 新文档标题
     * @param new_file 新文件
     *
     * @return HTTP响应
     */
    public HttpResponse<String> update_document(String datasetId, String documentId, String new_title, File new_file) throws Exception {

        fileTextExtractor = new FileTextExtractor();
        return Unirest.post("https://api.dify.ai/v1/datasets/"+datasetId+"/documents/"+documentId+"/update-by-text")
                .header("Authorization", Authorization_dataset)
                .header("Content-Type", "application/json")
                .body("{\n  \"name\": \""+new_title+"\",\n  \"text\": \""+fileTextExtractor.getExtractedText(new_file)+"\",\n  \"process_rule\": {\n    \"mode\": \"automatic\",\n    \"rules\": {\n      \"pre_processing_rules\": [\n        {\n          \"id\": \"remove_extra_spaces\",\n          \"enabled\": true\n        }\n      ],\n      \"segmentation\": {\n        \"separator\": \"###\",\n        \"max_tokens\": 500\n      }\n    }\n  }\n}")
                .asString();
    }

    /**
     * 删除文档
     * @param datasetId 数据集ID
     * @param documentId 文档ID
     * @return HTTP响应
     */
    public HttpResponse<String> delete_document(String datasetId, String documentId) {
        return Unirest.delete("https://api.dify.ai/v1/datasets/"+datasetId+"/documents/"+documentId)
                .header("Authorization", Authorization_dataset)
                .asString();
    }


    /**
     * 与大模型对话
     * @param query 用户查询内容
     * @return HTTP响应
     */
    public HttpResponse<String> chatLLM(String query) {
        return Unirest.post("https://api.dify.ai/v1/chat-messages")
                .header("Authorization", Authorization_app)
                .header("Content-Type", "application/json")
                .body("{\n  \"inputs\": {},\n  \"response_mode\": \"streaming\",\n  \"auto_generate_name\": true,\n  \"query\": \""+query+"\",\n  \"user\": \""+user_id+"\"\n}")
                .asString();
    }


}