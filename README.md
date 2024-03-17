
<img src="docs/images/Cake.svg"/>

## Component definition

- **Rule**: A series of condition define by promotion teams. If the conditions match, the **effect** will be executed. In simple language, this is **if statement**. Eg: "If this user is in 100 first log in users and user top-up the mobile phone".
- **Effect**: The **effect** after user finish an action that met  **rule** dragged on it. This is the **then** statement in **if-then** clause. Eg: "discount 30%".
- **Rule-engine**: Solve the rule-effect relationship by do **rule** validation, then do the **effect**.
- **Offers**: Store the rule-effect metadata and pass needed input the the **rule-engine**
- **Campaigns**: Define by promotion team, the campaign contain offers or in other word, the offers will be appeared when the campaign show it.
- **Invalid cache**: Store the invalidate user-campaign-offer metadata. The 
- **user-campaign-offer data**: The composite key define the association from user data, campaign and offer metadata. This can be use for some meaning: The user can or cannot use the offer from the campaign.
- **Validate campaign**: Do validation of  campaign-offer key. This can be use to update is campaign is still activating in realtime.
- *Validation*: Do validation on user-campaign-offer.
- **User Event**: An pipeline queue that received and keep event from user, including payment action.
- *Payment*: Do necessary step when user finisher paying action. This's subscribe to UserEvent queue to get the input data.
- *User*: Service to handle user information.
- **Pro-processed database**: The data will be pre-calculator first and try to push them to this database. This will improve time for validation and suggest offers for users. This should be a READ optimization database
- **Data processing**:  Handle the data pre-calculation and push the centralize database.
- **User's balance**: Do the money balance for user.



<img src="docs/images/Cake-get-offers.svg">

<img src="docs/images/Cake-campaign-creator.svg">

<img src="docs/images/Cake-issue-offer.svg">
