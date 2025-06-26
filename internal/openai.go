package internal

import (
	"bytes"
	"fmt"
	"log"

	"encoding/json"
	"text/template"

	"github.com/go-resty/resty/v2"
)

type InstructionContext struct {
	Json string
}

func OpenAIEval(apiKey string, ctx Context, command string) []Patch {
	ctxJson, _ := json.MarshalIndent(ctx, "", "  ")

	instructionCtx := InstructionContext{
		Json: string(ctxJson),
	}

	if apiKey == "" {
		log.Fatal("OpenAI API key not set")
	}
	
	var result []Patch

	endpoint := "https://api.openai.com/v1/chat/completions"

	client := resty.New()

	instructionTmplStr := `
You are a helpful virtual assistant that helps developers transform their code. You are given a code snippet below and are asked to make the transformation in the message. If the user's message is clear, call the get_patch_operations function.

The context is in JSON and is between the "#####" lines.

#####
{{ .Json }}
#####
`
	tmpl, err := template.New("instructionTemplate").Parse(instructionTmplStr)
	if err != nil {
		panic(err)
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, instructionCtx)
	if err != nil {
		panic(err)
	}

	instructions := outputBuffer.String()

	fmt.Println("instructions", instructions)

	headers := map[string]string{
		"Content-Type": "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": instructions,
			},
			{
				"role": "user",
				"content": command,
			},
		},
		"functions": []map[string]interface{}{
			{
				"name": "get_patch_operations",
				"description": "Get the list of patch operations (content, file, or directory level)",
				"parameters": map[string]interface{}{
					"type":    "object",
					"properties": map[string]interface{}{
						"patches": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"oneOf": []map[string]interface{}{
									{
										"type": "object",
										"properties": map[string]interface{}{
											"File": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"FileOp":      map[string]interface{}{"type": "string", "enum": []string{"FileCreate", "FileDelete", "FileRename", "DirectoryDelete", "DirectoryRename"}},
													"SourcePath":  map[string]interface{}{"type": "string"},
													"TargetPath":  map[string]interface{}{"type": "string"},
													"Lines":       map[string]interface{}{
														"type": "array",
														"items": map[string]interface{}{
															"type": "string",
														},
													},
												},
												"required":           []string{"FileOp", "SourcePath", "TargetPath"},
												"additionalProperties": false,
											},
										},
										"required": []string{"File"},
										"not":     map[string]interface{}{"required": []string{"Content"}},
									},
									{
										"properties": map[string]interface{}{
											"Content": map[string]interface{}{
												"type": "object",
												"properties": map[string]interface{}{
													"PatchType":     map[string]interface{}{"type": "string", "enum": []string{"PatchReplace", "PatchDelete", "PatchInsert"}},
													"OldLineNumber": map[string]interface{}{"type": "integer"},
													"OldLineCount":  map[string]interface{}{"type": "integer"},
													"NewContent":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
												},
												"required":           []string{"PatchType", "OldLineNumber", "OldLineCount", "NewContent"},
												"additionalProperties": false,
											},
										},
										"required": []string{"Content"},
										"not":     map[string]interface{}{"required": []string{"File"}},
									},
								},
							},
						},
					},
					"required": []string{"patches"},
					"additionalProperties": false,
				},
			},
		},
	}

	json, _ := json.MarshalIndent(data, "", "  ")

	fmt.Println("request", string(json))

	resp, _ := client.R().
		SetHeaders(headers).
		SetBody(json).
		Post(endpoint)

	fmt.Println("response", resp.String())
		
	return result
}