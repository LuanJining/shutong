package com.luanjining.tool;

import kong.unirest.HttpResponse;
import kong.unirest.Unirest;

public class QueryTool {

    private static final String Authorization_dataset = "Bearer dataset-juR2DRl8alE1eBvQmvJ6KiVs";

    /**
     * 查询文档
     * 若keyword为空则查询全部，默认分页，每页20条
     * @param keyword 关键词
     * @return HTTP响应
     */
    public HttpResponse<String> query_document(String keyword) {
        return  Unirest.get("https://api.dify.ai/v1/datasets/bafabbc4-0093-49ec-94a7-71a1e786105d/documents?page=1&limit=20&keyword="+keyword)
                .header("Authorization", Authorization_dataset)
                .asString();
    }
}
