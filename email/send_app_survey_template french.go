package sendemail

const AppSurveyTemplateFrench = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Subject}}</title>
	
	<style type="text/css">
	@import url('https://fonts.googleapis.com/css2?family=Afacad:wght@400;500;600;700&display=swap'); /* Added Afacad */
	@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
	
	</style>
</head>
<body style="margin: 0px;">

	<center class="wrapper" style="width: 100%; table-layout: fixed; background-color: #DAD4CB; padding-bottom: 60px;">

		<table class="main" style="border-spacing: 0; width: 100%; max-width: 666px; background-color: #FDFBF6; font-family: 'Inter', Arial, Helvetica, sans-serif; color: #72513C; text-align: center;" width="100%">

		<tr class="logo-wrapper" style="background-color: #72513C; height: 66px;">
			<td style=""">
				<a href="https://www.ommindshop.com" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/logo.png" alt="Ommind" width="147" height="34" style=""></a>
			</td>
		</tr>
					
		<tr > 
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 500px;">
					<tr>
						<td align="center">
							<p style="font-weight: 600px; font-size: 30px;">Bonjour à toi</p>
							<p style="font-size: 14px;">Merci encore d’avoir passé notre test des chakras — nous espérons qu’il t’a apporté plus de clarté sur l’endroit où tu te trouves dans ton cheminement intérieur.</p>
							<p style="font-size: 14px; margin-bottom: 0px">Nous sommes en train de créer quelque chose de spécial :</p>
							<p style="font-size: 14px; margin-top: 0px">Une application de méditation et de croissance spirituelle personnalisée, inspirée de la sagesse tibétaine, conçue pour te guider pas à pas — du calme, à la transformation, jusqu’à l’éveil.</p>
							<p style="font-size: 14px; margin-bottom: 0px;">Nous aimerions beaucoup ton aide.</p>
							<p style="font-size: 14px; margin-bottom: 0px; margin-top: 0px;"">Aurais-tu <strong>2 à 3 minutes</strong> pour répondre à un petit questionnaire qui nous aidera à façonner la prochaine version d’<strong>OmMind</strong> ?</p>
							<p style="font-size: 14px; margin-bottom: 0px">En remerciement :</p>
							<p style="font-size: 14px; margin-bottom: 0px; margin-top: 0px;">✨ 1 personne sur 5 recevra un <strong>bracelet chakra en pierres naturelles</strong></p>
							<p style="font-size: 14px; margin-bottom: 0px; margin-top: 0px;">💎 Et tout le monde recevra <strong>-50 % sur la boutique OmMind</strong></p>
							<p style="font-size: 14px;">Ta voix compte vraiment — que tu débutes ou que tu sois déjà engagé(e) sur ce chemin depuis longtemps.</p>
							<p style="font-size: 14px"><a target="_blank" href="https://docs.google.com/forms/d/e/1FAIpQLScy75LVfzwo_NJa-qAXDHaFL_vusNppajUl7Lg_VFy_6OphcQ/viewform">[ Commencer le questionnaire ]</a></p>
						</td>
					</tr>
					
				</table>
			</td>
		</tr>


		<tr > 
			<td>
				<table align="center" role="presentation" width="100%" style="max-width: 600px; margin-bottom: 20px;">
					<tr>
						<td align="center">
							<tr>
								<td height="234" width="601" align="center" valign="middle" background="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/xr.png" style="height: 234px; width: 601px; background-image: url('https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/xr.png'); background-repeat: no-repeat; background-position: center center; background-size: cover; text-align: center; vertical-align: middle;">
									<div>
										<table align="center" role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 500px;"> 
											<tr>
												<td align="center" style="padding: 20px 10px;"> 													
												</td>
											</tr>
										</table>
									</div>
								</td>
							</tr>
						</td>
					</tr>
					
				</table>
			</td>
		</tr>

		<tr>
			<td style="padding:0">
				<table width="100%" style="background-color: #72513C; border-spacing: 0px;">
					<tr>
						<td>
							<tr class="logo-wrapper">
								<td align="center" style="padding-top: 20px;">
									<a href="https://www.ommindshop.com/" target="_blank"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/logo.png" alt="Ommind" width="147" height="34" style=""></a>
								</td>
							</tr>						
						</td>
					</tr>
					<tr >
						<td align="center" style="padding: 25px 25px; font-size: 0;">
							<a href="https://www.facebook.com/profile.php?id=61573346513645" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-1.png"  alt="Social Media 1" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.tiktok.com/@ommindshop" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-2.png"  alt="Social Media 2" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.instagram.com/ommind_shop/" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-3.png"  alt="Social Media 3" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
							<a href="https://www.youtube.com/@OmMind-Official" target="_blank" style="text-decoration: none; display: inline-block; vertical-align: middle;"><img src="https://ommind-public.s3.eu-west-2.amazonaws.com/chakara-report-email-tempalte-images/company-info-logo-4.png"  alt="Social Media 4" width="25" height="25" style="width: 25px; height: 25px; border: 0; display: inline-block; margin: 0 12px; vertical-align: middle;"></a>
						</td>
					</tr>
					<tr>
						<td align="center" style="padding: 0 25px 50px 25px;"> 
							<p style="font-size: 14px; color: white; margin-bottom: 5px; margin-top: 0; margin-bottom: 10px;">86-90 Paul Street, London, EC2A 4UX, United Kingdom</p>
							<p style="font-size: 14px; color: #BDBDBD; margin-top: 0; margin-bottom: 0;">COPYRIGHT ©2025 OmMind Shop All Rights Reserved</p>
						</td>
					</tr>
					
					
				</table>
			</td>
		</tr>			
		</table> 

	</center> 
</body>
</html>
`
