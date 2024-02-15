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
    const noOfStudents = parseInt(form.elements['captone_noStudent'].value);
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

    const curl = `http://localhost:5002/api/v1/records`;
    console.log(curl);

    request.open("POST", curl);

    request.send(JSON.stringify ({
        "name": name,
        "roleOfContact": roleOfContact,
        "noOfStudents": noOfStudents,
        "acadYr": acadYr,
        "capstoneTitle": capstoneTitle,
        "companyName": companyName,
        "companyContact": companyContact,
        "projDesc": projDesc,
    }));

    console.log("Request Status: ", request.status);
    alert("Capstone record is created.")
    form.reset();

}

function listCapstones() {
    const url = `http://localhost:5002/api/v1/records/all`;
    fetch(url)
        .then(response => {
            if(!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then( data => {
            console.log("Data from server: ", data);

            var tableBody = document.getElementById('allcapstone');

            tableBody.innerHTML = '';

            data.forEach(record => {
                var row = tableBody.insertRow();
                //Edit the Javascript Function !!!!
                row.innerHTML = `   <th scope="row">${record.recordId}</th>
                                    <td>${record.name}</td>
                                    <td>${record.capstoneTitle}</td>
                                    <td>${record.roleOfContact}</td>
                                    <td>${record.noOfStudents}</td>
                                    <td>${record.companyName}</td>
                                    <td>${record.companyContact}</td>
                                    <td>${record.acadYr}</td>
                                    <td class="description">${record.projDesc}</td>
                                    <td>
                                        <button class="btn btn-outline-secondary" onclick="return modifyCapstone(${record.recordId})">Modify</button>
                                        <button class="btn btn-outline-secondary" onclick="return deleteCapstone(${record.recordId})">Delete</button>
                                    </td>
                                `;
            });
        })
        .catch(error => console.error('Error fetching record details: ', error))
}

function queryCapstone() {
    var request = new XMLHttpRequest();
    const form = document.getElementById('querycapstone');

    const curl = 'http://localhost:5001/api/v1/records/search';

    //HTML VALUE 
    const queryacadYr = form.elements['query_acadYr'].value;
    const queryKeyword = form.elements['query_keyword'].value;

    //CHECK
    console.log(queryacadYr);
    console.log(queryKeyword);

    
}
