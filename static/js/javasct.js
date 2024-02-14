/*

    router.HandleFunc("/api/v1/records/all", ListAllRecordsHandler).Methods("GET")
	router.HandleFunc("/api/v1/records", CreateRecordHandler).Methods("POST")
	router.HandleFunc("/api/v1/records/delete", DeleteRecordHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/records/{recordID}", UpdateRecordHandler).Methods("PUT")
	router.HandleFunc("/api/v1/records/search", QueryRecordHandler).Methods("GET")


type Record struct {
	RecordID       int    `json:"recordId"`
	Name           string `json:"name"`
	RoleOfContact  string `json:"roleOfContact"`
	NoOfStudents   int    `json:"noOfStudents"`
	AcadYr         string `json:"acadYr"`
	CapstoneTitle  string `json:"capstoneTitle"`
	CompanyName    string `json:"companyName"`
	CompanyContact string `json:"companyContact"`
	ProjDesc       string `json:"projDesc"`
}
 */

function createCapstone() {
    var request = new XMLHttpRequest();
    const form = document.getElementById('createcapstoneForm');

    //FORM VALUES
    const name = form.elements['captone_name'].value;
    const roleOfContact = document.querySelector('input[name="flexRadioDefault"]:checked').value;
    const noOfStudents = form.elements['captone_noStudent'].value;
    const acadYr = form.elements['captone_academicYear'].value;
    const capstoneTitle = form.elements['captone_capstonetitle'].value;
    const companyName = form.elements['captone_companyName'].value;
    const companyContact = form.elements['captone_companyPOC'].value;
    const projDesc = form.elements['captone_description'].value;

    //CONSOLE LOG CHECK
    console.log(name);
    console.log(roleOfContact);
    console.log(noOfStudents);
    console.log(acadYr);
    console.log(capstoneTitle);
    console.log(companyName);
    console.log(companyName);
    console.log(companyContact);
    console.log(projDesc);

    const curl = 'http://localhost:5001/api/v1/records';
    console.log(curl);

    request.open("POST", curl);

    const newcapstoneJSON = JSON.stringify ({
        "name": name,
        "roleOfContact": roleOfContact,
        "noOfStudents": noOfStudents,
        "acadYr": acadYr,
        "capstoneTitle": capstoneTitle,
        "companyName": companyName,
        "companyContact": companyContact,
        "projDesc": projDesc,
    });

    console.log(newcapstoneJSON);

    request.send(newcapstoneJSON);

    const reqStatus = request.status;
    console.log(reqStatus)
    
    if(reqStatus == 0) {
        alert("Capstone record is not created.");
    } else if (reqStatus == 200) {
        alert("Capstone record created successfully.")
    } else if (reqStatus == 400) {
        const errorData = JSON.parse(request.responseText);
        alert("Error creating capstone: " + errorData.message);
    } else if (reqStatus == 404) {
        alert("Endpoint not found.");
    } else {
        alert("An error occured, please try again later.");
    }

    form.reset();

}
