import com.luanjining.API_demo;
import kong.unirest.HttpResponse;
import kong.unirest.Unirest;

import java.io.File;

public class Test {

    public static void main(String[] args) throws Exception {
        API_demo api = new API_demo();



        //create_dataset测试
//        System.out.printf(api.create_dataset("test1","test").getBody());

        //create_document测试
//        System.out.printf(api.create_document("bafabbc4-0093-49ec-94a7-71a1e786105d","test",new File("Luan_RAG/src/test/java/test.txt")).getBody());

        //update_document测试
//        System.out.printf(api.update_document("bafabbc4-0093-49ec-94a7-71a1e786105d","991f797e-bed4-4b66-86e4-c76edf97fe4b","test_update",new File("Luan_RAG/src/test/java/test.txt")).getBody());

        //delete_document测试
//        System.out.printf(api.delete_document("bafabbc4-0093-49ec-94a7-71a1e786105d","31323f73-dc15-4a6a-9bbc-f0fc8bc4d855").getBody());

        //chatLLM测试
        System.out.printf(api.chatLLM("你好，你是谁").getBody());


    }


}
