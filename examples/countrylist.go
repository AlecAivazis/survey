package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var countryQs = []*survey.Question{
	{
		Name: "country",
		Prompt: survey.NewSingleSelect().SetMessage("Choose a country:").
				AddOption("Afghanistan", nil, false).
				AddOption("Åland Islands", nil, false).
				AddOption("Albania", nil, false).
				AddOption("Algeria", nil, false).
				AddOption("American Samoa", nil, false).
				AddOption("AndorrA", nil, false).
				AddOption("Angola", nil, false).
				AddOption("Anguilla", nil, false).
				AddOption("Antarctica", nil, false).
				AddOption("Antigua and Barbuda", nil, false).
				AddOption("Argentina", nil, false).
				AddOption("Armenia", nil, false).
				AddOption("Aruba", nil, false).
				AddOption("Australia", nil, false).
				AddOption("Austria", nil, false).
				AddOption("Azerbaijan", nil, false).
				AddOption("Bahamas", nil, false).
				AddOption("Bahrain", nil, false).
				AddOption("Bangladesh", nil, false).
				AddOption("Barbados", nil, false).
				AddOption("Belarus", nil, false).
				AddOption("Belgium", nil, false).
				AddOption("Belize", nil, false).
				AddOption("Benin", nil, false).
				AddOption("Bermuda", nil, false).
				AddOption("Bhutan", nil, false).
				AddOption("Bolivia", nil, false).
				AddOption("Bosnia and Herzegovina", nil, false).
				AddOption("Botswana", nil, false).
				AddOption("Bouvet Island", nil, false).
				AddOption("Brazil", nil, false).
				AddOption("British Indian Ocean Territory", nil, false).
				AddOption("Brunei Darussalam", nil, false).
				AddOption("Bulgaria", nil, false).
				AddOption("Burkina Faso", nil, false).
				AddOption("Burundi", nil, false).
				AddOption("Cambodia", nil, false).
				AddOption("Cameroon", nil, false).
				AddOption("Canada", nil, false).
				AddOption("Cape Verde", nil, false).
				AddOption("Cayman Islands", nil, false).
				AddOption("Central African Republic", nil, false).
				AddOption("Chad", nil, false).
				AddOption("Chile", nil, false).
				AddOption("China", nil, false).
				AddOption("Christmas Island", nil, false).
				AddOption("Cocos (Keeling) Islands", nil, false).
				AddOption("Colombia", nil, false).
				AddOption("Comoros", nil, false).
				AddOption("Congo", nil, false).
				AddOption("Congo, The Democratic Republic of the", nil, false).
				AddOption("Cook Islands", nil, false).
				AddOption("Costa Rica", nil, false).
				AddOption("Cote D'Ivoire", nil, false).
				AddOption("Croatia", nil, false).
				AddOption("Cuba", nil, false).
				AddOption("Cyprus", nil, false).
				AddOption("Czech Republic", nil, false).
				AddOption("Denmark", nil, false).
				AddOption("Djibouti", nil, false).
				AddOption("Dominica", nil, false).
				AddOption("Dominican Republic", nil, false).
				AddOption("Ecuador", nil, false).
				AddOption("Egypt", nil, false).
				AddOption("El Salvador", nil, false).
				AddOption("Equatorial Guinea", nil, false).
				AddOption("Eritrea", nil, false).
				AddOption("Estonia", nil, false).
				AddOption("Ethiopia", nil, false).
				AddOption("Falkland Islands (Malvinas)", nil, false).
				AddOption("Faroe Islands", nil, false).
				AddOption("Fiji", nil, false).
				AddOption("Finland", nil, false).
				AddOption("France", nil, false).
				AddOption("French Guiana", nil, false).
				AddOption("French Polynesia", nil, false).
				AddOption("French Southern Territories", nil, false).
				AddOption("Gabon", nil, false).
				AddOption("Gambia", nil, false).
				AddOption("Georgia", nil, false).
				AddOption("Germany", nil, false).
				AddOption("Ghana", nil, false).
				AddOption("Gibraltar", nil, false).
				AddOption("Greece", nil, false).
				AddOption("Greenland", nil, false).
				AddOption("Grenada", nil, false).
				AddOption("Guadeloupe", nil, false).
				AddOption("Guam", nil, false).
				AddOption("Guatemala", nil, false).
				AddOption("Guernsey", nil, false).
				AddOption("Guinea", nil, false).
				AddOption("Guinea-Bissau", nil, false).
				AddOption("Guyana", nil, false).
				AddOption("Haiti", nil, false).
				AddOption("Heard Island and Mcdonald Islands", nil, false).
				AddOption("Holy See (Vatican City State)", nil, false).
				AddOption("Honduras", nil, false).
				AddOption("Hong Kong", nil, false).
				AddOption("Hungary", nil, false).
				AddOption("Iceland", nil, false).
				AddOption("India", nil, false).
				AddOption("Indonesia", nil, false).
				AddOption("Iran, Islamic Republic Of", nil, false).
				AddOption("Iraq", nil, false).
				AddOption("Ireland", nil, false).
				AddOption("Isle of Man", nil, false).
				AddOption("Israel", nil, false).
				AddOption("Italy", nil, false).
				AddOption("Jamaica", nil, false).
				AddOption("Japan", nil, false).
				AddOption("Jersey", nil, false).
				AddOption("Jordan", nil, false).
				AddOption("Kazakhstan", nil, false).
				AddOption("Kenya", nil, false).
				AddOption("Kiribati", nil, false).
				AddOption("Korea, Democratic People'S Republic of", nil, false).
				AddOption("Korea, Republic of", nil, false).
				AddOption("Kuwait", nil, false).
				AddOption("Kyrgyzstan", nil, false).
				AddOption("Lao People'S Democratic Republic", nil, false).
				AddOption("Latvia", nil, false).
				AddOption("Lebanon", nil, false).
				AddOption("Lesotho", nil, false).
				AddOption("Liberia", nil, false).
				AddOption("Libyan Arab Jamahiriya", nil, false).
				AddOption("Liechtenstein", nil, false).
				AddOption("Lithuania", nil, false).
				AddOption("Luxembourg", nil, false).
				AddOption("Macao", nil, false).
				AddOption("Macedonia, The Former Yugoslav Republic of", nil, false).
				AddOption("Madagascar", nil, false).
				AddOption("Malawi", nil, false).
				AddOption("Malaysia", nil, false).
				AddOption("Maldives", nil, false).
				AddOption("Mali", nil, false).
				AddOption("Malta", nil, false).
				AddOption("Marshall Islands", nil, false).
				AddOption("Martinique", nil, false).
				AddOption("Mauritania", nil, false).
				AddOption("Mauritius", nil, false).
				AddOption("Mayotte", nil, false).
				AddOption("Mexico", nil, false).
				AddOption("Micronesia, Federated States of", nil, false).
				AddOption("Moldova, Republic of", nil, false).
				AddOption("Monaco", nil, false).
				AddOption("Mongolia", nil, false).
				AddOption("Montserrat", nil, false).
				AddOption("Morocco", nil, false).
				AddOption("Mozambique", nil, false).
				AddOption("Myanmar", nil, false).
				AddOption("Namibia", nil, false).
				AddOption("Nauru", nil, false).
				AddOption("Nepal", nil, false).
				AddOption("Netherlands", nil, false).
				AddOption("Netherlands Antilles", nil, false).
				AddOption("New Caledonia", nil, false).
				AddOption("New Zealand", nil, false).
				AddOption("Nicaragua", nil, false).
				AddOption("Niger", nil, false).
				AddOption("Nigeria", nil, false).
				AddOption("Niue", nil, false).
				AddOption("Norfolk Island", nil, false).
				AddOption("Northern Mariana Islands", nil, false).
				AddOption("Norway", nil, false).
				AddOption("Oman", nil, false).
				AddOption("Pakistan", nil, false).
				AddOption("Palau", nil, false).
				AddOption("Palestinian Territory, Occupied", nil, false).
				AddOption("Panama", nil, false).
				AddOption("Papua New Guinea", nil, false).
				AddOption("Paraguay", nil, false).
				AddOption("Peru", nil, false).
				AddOption("Philippines", nil, false).
				AddOption("Pitcairn", nil, false).
				AddOption("Poland", nil, false).
				AddOption("Portugal", nil, false).
				AddOption("Puerto Rico", nil, false).
				AddOption("Qatar", nil, false).
				AddOption("Reunion", nil, false).
				AddOption("Romania", nil, false).
				AddOption("Russian Federation", nil, false).
				AddOption("RWANDA", nil, false).
				AddOption("Saint Helena", nil, false).
				AddOption("Saint Kitts and Nevis", nil, false).
				AddOption("Saint Lucia", nil, false).
				AddOption("Saint Pierre and Miquelon", nil, false).
				AddOption("Saint Vincent and the Grenadines", nil, false).
				AddOption("Samoa", nil, false).
				AddOption("San Marino", nil, false).
				AddOption("Sao Tome and Principe", nil, false).
				AddOption("Saudi Arabia", nil, false).
				AddOption("Senegal", nil, false).
				AddOption("Serbia and Montenegro", nil, false).
				AddOption("Seychelles", nil, false).
				AddOption("Sierra Leone", nil, false).
				AddOption("Singapore", nil, false).
				AddOption("Slovakia", nil, false).
				AddOption("Slovenia", nil, false).
				AddOption("Solomon Islands", nil, false).
				AddOption("Somalia", nil, false).
				AddOption("South Africa", nil, false).
				AddOption("South Georgia and the South Sandwich Islands", nil, false).
				AddOption("Spain", nil, false).
				AddOption("Sri Lanka", nil, false).
				AddOption("Sudan", nil, false).
				AddOption("Suriname", nil, false).
				AddOption("Svalbard and Jan Mayen", nil, false).
				AddOption("Swaziland", nil, false).
				AddOption("Sweden", nil, false).
				AddOption("Switzerland", nil, false).
				AddOption("Syrian Arab Republic", nil, false).
				AddOption("Taiwan, Province of China", nil, false).
				AddOption("Tajikistan", nil, false).
				AddOption("Tanzania, United Republic of", nil, false).
				AddOption("Thailand", nil, false).
				AddOption("Timor-Leste", nil, false).
				AddOption("Togo", nil, false).
				AddOption("Tokelau", nil, false).
				AddOption("Tonga", nil, false).
				AddOption("Trinidad and Tobago", nil, false).
				AddOption("Tunisia", nil, false).
				AddOption("Turkey", nil, false).
				AddOption("Turkmenistan", nil, false).
				AddOption("Turks and Caicos Islands", nil, false).
				AddOption("Tuvalu", nil, false).
				AddOption("Uganda", nil, false).
				AddOption("Ukraine", nil, false).
				AddOption("United Arab Emirates", nil, false).
				AddOption("United Kingdom", nil, false).
				AddOption("United States", nil, false).
				AddOption("United States Minor Outlying Islands", nil, false).
				AddOption("Uruguay", nil, false).
				AddOption("Uzbekistan", nil, false).
				AddOption("Vanuatu", nil, false).
				AddOption("Venezuela", nil, false).
				AddOption("Viet Nam", nil, false).
				AddOption("Virgin Islands, British", nil, false).
				AddOption("Virgin Islands, U.S.", nil, false).
				AddOption("Wallis and Futuna", nil, false).
				AddOption("Western Sahara", nil, false).
				AddOption("Yemen", nil, false).
				AddOption("Zambia", nil, false).
				AddOption("Zimbabwe", nil, false),
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		Country *survey.Option
	}{}

	// ask the question
	err := survey.Ask(countryQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("you chose %s.\n", answers.Country)
}
